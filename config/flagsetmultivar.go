package config


import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type FlagSetMultiVar struct {
	flag.FlagSet
	multis map[string][]string
	alternatives *[]string
	types map[string]string
}

func NewFlagSet(desc string) FlagSetMultiVar {
	fs := flag.NewFlagSet(desc, flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)
	return FlagSetMultiVar{*fs, map[string][]string{}, &[]string{}, map[string]string{}}
}


func (flagSet FlagSetMultiVar) MultiIntVar(p *int, items []string, def int, description string) {
	flagSet.Record(items, "int")
	for _, item := range items {
		flagSet.IntVar(p, item, def, description)
	}
}
func (flagSet FlagSetMultiVar) MultiStringVar(p *string, items []string, def string, description string) {
	flagSet.Record(items, "string")
	for _, item := range items {
		flagSet.StringVar(p, item, def, description)
	}
}
func (flagSet FlagSetMultiVar) MultiBoolVar(p *bool, items []string, def bool, description string) {
	flagSet.Record(items, "bool")
	for _, item := range items {
		flagSet.BoolVar(p, item, def, description)
	}
}
func (flagSet FlagSetMultiVar) isAlternative(name string) bool {
	for _, alt := range *flagSet.alternatives {
		if alt == name {
			return true
		}
	}
	return false
}

func (flagSet FlagSetMultiVar) Record(items []string, typ string) {
	var key string
	for i, item := range items {
		if i == 0 {
			key = item
			if _, ok := flagSet.multis[key]; !ok {
				flagSet.multis[key] = []string{}
			}
		} else {
			//alternative names ...
			//key is same as before
			*flagSet.alternatives = append(*flagSet.alternatives, item)
			flagSet.multis[key] = append(flagSet.multis[key], item)
		}
	}
	flagSet.types[key] = typ
	
}
func (flagSet FlagSetMultiVar) isBool(name string) bool {
	typ, ok := flagSet.types[name]
	if ok {
		if typ == "bool" {
			return true
		}
	}
	return false
}
func (flagSet FlagSetMultiVar) PrintDefaults() {
	flagSet.VisitAll(func(fl *flag.Flag) {
		format := "-%s=%s"
		l := 0
		alts, isMulti := flagSet.multis[fl.Name]
		if isMulti {
			li, _ := fmt.Fprintf(os.Stderr, "  ")
			l += li
			for _, alt := range alts {
				if len(alt)>1 {
					li, _ := fmt.Print("-")
					l += li
				}
				li, _ := fmt.Fprintf(os.Stderr, "-%s  ", alt)
				l += li
			}
			if len(fl.Name)>1 {
				li, _ := fmt.Print("-")
				l += li
			}
			if flagSet.isBool(fl.Name) {
				li, _ = fmt.Fprintf(os.Stderr, "-%s", fl.Name)
				l += li
			} else {
				li, _ = fmt.Fprintf(os.Stderr, format, fl.Name, fl.DefValue)
				l += li
			}
		} else {
			if !flagSet.isAlternative(fl.Name) {
				li, _ := fmt.Fprintf(os.Stderr, "  ")
				l += li
				if len(fl.Name)>1 {
					li, _ := fmt.Print("-")
					l += li
				}
				li, _ = fmt.Fprintf(os.Stderr, format, fl.Name, fl.DefValue)
				l += li
			}
		}
		if !flagSet.isAlternative(fl.Name) {
			for l < 25 {
				l += 1
				fmt.Fprintf(os.Stderr, " ")
			}
			fmt.Fprintf(os.Stderr, ": %s\n", fl.Usage)
		}
		
	})
	fmt.Println("")
}