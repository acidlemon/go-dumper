package dump

import (
	"fmt"
	r "reflect"
	"io"
	"os"
)

var emptyString = ""

// Prints to the writer the value with indentation.
func Fdump(out io.Writer, v_ interface{}) {
	// forward decl
	var dump0 func(r.Value, int)
	var dump func(r.Value, int, *string, *string)

	done := make(map[string]bool)

	dump = func(v r.Value, d int, prefix *string, suffix *string) {
		pad := func() {
			res := ""
			for i := 0; i < d; i++ {
				res += "  "
			}
			fmt.Fprintf(out, res)
		}

		padprefix := func() {
			if prefix != nil {
				fmt.Fprintf(out, *prefix)
			} else {
				res := ""
				for i := 0; i < d; i++ {
					res += "  "
				}
				fmt.Fprintf(out, res)
			}
		}

		printf := func(s string, args ...interface{}) { fmt.Fprintf(out, s, args...) }
		print := func(args ...interface{}) { fmt.Fprint(out, args...) }

		// prevent circular for composite types
		if v.CanAddr() {
			addr := v.Addr()
			key := fmt.Sprintf("0x%08x %v", addr.Pointer(), v.Type())
			if _, exists := done[key]; exists {
				padprefix()
				printf("<%s>", key)
				return
			} else {
				done[key] = true
			}
		}

		switch v.Kind() {
		case r.Array:
			padprefix()
			if v.IsNil() {
				printf("nil ([%v]%v)", v.Len(), v.Type().Elem())
			} else {
				printf("[%v]%v {\n", v.Len(), v.Type().Elem())
				for i := 0; i < v.Len(); i++ {
					dump0(v.Field(i), d+1)
					if i != v.Len()-1 {
						printf(",\n")
					} else {
						print("\n")
					}
				}
				pad()
				print("}")
			}

		case r.Slice:
			padprefix()
			if v.IsNil() {
				printf("nil ([]%v)", v.Type().Elem())
			} else {
				printf("[]%v (len=%d) {\n", v.Type().Elem(), v.Len())
				for i := 0; i < v.Len(); i++ {
					dump0(v.Index(i), d+1)
					if i != v.Len()-1 {
						printf(",\n")
					} else {
						print("\n")
					}
				}
				pad()
				print("}")
			}

		case r.Map:
			padprefix()
			t := v.Type()
			if v.IsNil() {
				printf("nil (map[%s]%v)", t.Key().Name(), t.Elem())
			} else {
				printf("map[%s]%v {\n", t.Key().Name(), t.Elem())
				for i, k := range v.MapKeys() {
					dump0(k, d+1)
					printf(": ")
					dump(v.MapIndex(k), d+1, &emptyString, nil)
					if i != v.Len()-1 {
						printf(",\n")
					} else {
						print("\n")
					}
				}
				pad()
				print("}")
			}

		case r.Ptr:
			padprefix()
			if v.IsNil() {
				printf("nil (*%s)", v.Type().Elem())
			} else {
				printf("(0x%08x) &", v.Pointer())
				dump(v.Elem(), d, &emptyString, nil)
			}

		case r.Struct:
			padprefix()
			t := v.Type()
			printf("%v {\n", t.Name())
			d += 1
			for i := 0; i < v.NumField(); i++ {
				pad()
				printf("%v", t.Field(i).Name)
				printf(": ")
				dump(v.Field(i), d, &emptyString, nil)
				if i != v.NumField()-1 {
					printf(",\n")
				} else {
					print("\n")
				}
			}
			d -= 1
			pad()
			print("}")

		case r.Interface:
			padprefix()
			t := v.Type()
			printf("interface(%s) ", t.Name())
			//dump(v.Elem(), d, &emptyString, nil)

		case r.String:
			padprefix()
			printf("\"%v\"", v.String())

		case r.Bool:
			padprefix()
			printf("%v", v.Bool())

		case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
			padprefix()
			printf("%v (%s)", v.Int(), v.Type().Name())

		case r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64:
			padprefix()
			printf("%v (%s)", v.Uint(), v.Type().Name())

		case r.Float32, r.Float64:
			padprefix()
			printf("%v (%s)", v.Float(), v.Type().Name())

		case r.Complex64, r.Complex128:
			padprefix()
			printf("%v (%s)", v.Complex(), v.Type().Name())

		case r.Uintptr:
			padprefix()
			printf("0x%x (%s)", v.Uint(), v.Type().Name())

		case r.UnsafePointer:
			padprefix()
			printf("0x%x (unsafe %s)", v.Pointer(), v.Type().Name())

		// Chan & Func
		default:
			padprefix()
			if v.IsValid() {
				if v.IsNil() {
					printf("nil")
				} else {
					printf("%v (%v)", v.Type(), v.Type().Name())
				}
			}
		}
	}

	dump0 = func(v r.Value, d int) { dump(v, d, nil, nil) }

	v := r.ValueOf(v_)
	dump0(v, 0)
	fmt.Fprintf(out, "\n")
}

// Print to standard out the value that is passed as the argument with indentation.
// Pointers are dereferenced.
func Dump(v_ interface{}) { Fdump(os.Stdout, v_) }
