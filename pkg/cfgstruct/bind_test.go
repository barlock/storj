// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package cfgstruct

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/pflag"
)

func assertEqual(actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		panic(fmt.Sprintf("expected %v, got %v", expected, actual))
	}
}

func TestBind(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.PanicOnError)
	var c struct {
		String   string        `default:""`
		Bool     bool          `releaseDefault:"false" devDefault:"true"`
		Int64    int64         `releaseDefault:"0" devDefault:"1"`
		Int      int           `default:"0"`
		Uint64   uint64        `default:"0"`
		Uint     uint          `default:"0"`
		Float64  float64       `default:"0"`
		Duration time.Duration `default:"0"`
		Struct   struct {
			AnotherString string `default:""`
		}
		Fields [10]struct {
			AnotherInt int `default:"0"`
		}
	}
	Bind(f, &c, UseReleaseDefaults())

	assertEqual(c.String, string(""))
	assertEqual(c.Bool, bool(false))
	assertEqual(c.Int64, int64(0))
	assertEqual(c.Int, int(0))
	assertEqual(c.Uint64, uint64(0))
	assertEqual(c.Uint, uint(0))
	assertEqual(c.Float64, float64(0))
	assertEqual(c.Duration, time.Duration(0))
	assertEqual(c.Struct.AnotherString, string(""))
	assertEqual(c.Fields[0].AnotherInt, int(0))
	assertEqual(c.Fields[3].AnotherInt, int(0))
	err := f.Parse([]string{
		"--string=1",
		"--bool=true",
		"--int64=1",
		"--int=1",
		"--uint64=1",
		"--uint=1",
		"--float64=1",
		"--duration=1h",
		"--struct.another-string=1",
		"--fields.03.another-int=1"})
	if err != nil {
		panic(err)
	}
	assertEqual(c.String, string("1"))
	assertEqual(c.Bool, bool(true))
	assertEqual(c.Int64, int64(1))
	assertEqual(c.Int, int(1))
	assertEqual(c.Uint64, uint64(1))
	assertEqual(c.Uint, uint(1))
	assertEqual(c.Float64, float64(1))
	assertEqual(c.Duration, time.Hour)
	assertEqual(c.Struct.AnotherString, string("1"))
	assertEqual(c.Fields[0].AnotherInt, int(0))
	assertEqual(c.Fields[3].AnotherInt, int(1))
}

func TestConfDir(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.PanicOnError)
	var c struct {
		String    string `default:"-$CONFDIR+"`
		MyStruct1 struct {
			String    string `default:"1${CONFDIR}2"`
			MyStruct2 struct {
				String string `default:"2${CONFDIR}3"`
			}
		}
	}
	Bind(f, &c, UseReleaseDefaults(), ConfDir("confpath"))
	assertEqual(f.Lookup("string").DefValue, "-confpath+")
	assertEqual(f.Lookup("my-struct1.string").DefValue, "1confpath2")
	assertEqual(f.Lookup("my-struct1.my-struct2.string").DefValue, "2confpath3")
}

func TestNesting(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.PanicOnError)
	var c struct {
		String    string `default:"-$CONFDIR+"`
		MyStruct1 struct {
			String    string `default:"1${CONFDIR}2"`
			MyStruct2 struct {
				String string `default:"2${CONFDIR}3"`
			}
		}
	}
	Bind(f, &c, UseReleaseDefaults(), ConfDirNested("confpath"))
	assertEqual(f.Lookup("string").DefValue, "-confpath+")
	assertEqual(f.Lookup("my-struct1.string").DefValue, filepath.FromSlash("1confpath/my-struct12"))
	assertEqual(f.Lookup("my-struct1.my-struct2.string").DefValue, filepath.FromSlash("2confpath/my-struct1/my-struct23"))
}

func TestBindDevDefaults(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.PanicOnError)
	var c struct {
		String   string        `default:"dev"`
		Bool     bool          `releaseDefault:"false" devDefault:"true"`
		Int64    int64         `releaseDefault:"0" devDefault:"1"`
		Int      int           `default:"2"`
		Uint64   uint64        `default:"3"`
		Uint     uint          `releaseDefault:"0" devDefault:"4"`
		Float64  float64       `default:"5.5"`
		Duration time.Duration `default:"1h"`
		Struct   struct {
			AnotherString string `default:"dev2"`
		}
		Fields [10]struct {
			AnotherInt int `default:"6"`
		}
	}
	Bind(f, &c, UseDevDefaults())

	assertEqual(c.String, string("dev"))
	assertEqual(c.Bool, bool(true))
	assertEqual(c.Int64, int64(1))
	assertEqual(c.Int, int(2))
	assertEqual(c.Uint64, uint64(3))
	assertEqual(c.Uint, uint(4))
	assertEqual(c.Float64, float64(5.5))
	assertEqual(c.Duration, time.Hour)
	assertEqual(c.Struct.AnotherString, string("dev2"))
	assertEqual(c.Fields[0].AnotherInt, int(6))
	assertEqual(c.Fields[3].AnotherInt, int(6))
	err := f.Parse([]string{
		"--string=1",
		"--bool=true",
		"--int64=1",
		"--int=1",
		"--uint64=1",
		"--uint=1",
		"--float64=1",
		"--duration=1h",
		"--struct.another-string=1",
		"--fields.03.another-int=1"})
	if err != nil {
		panic(err)
	}
	assertEqual(c.String, string("1"))
	assertEqual(c.Bool, bool(true))
	assertEqual(c.Int64, int64(1))
	assertEqual(c.Int, int(1))
	assertEqual(c.Uint64, uint64(1))
	assertEqual(c.Uint, uint(1))
	assertEqual(c.Float64, float64(1))
	assertEqual(c.Duration, time.Hour)
	assertEqual(c.Struct.AnotherString, string("1"))
	assertEqual(c.Fields[0].AnotherInt, int(6))
	assertEqual(c.Fields[3].AnotherInt, int(1))
}
