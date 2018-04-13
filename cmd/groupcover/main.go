// Copyright 2016 by Leipzig University Library, http://ub.uni-leipzig.de
//                    The Finc Authors, http://finc.info
//                    Martin Czygan, <martin.czygan@uni-leipzig.de>
//
// This file is part of some open source application.
//
// Some open source application is free software: you can redistribute
// it and/or modify it under the terms of the GNU General Public
// License as published by the Free Software Foundation, either
// version 3 of the License, or (at your option) any later version.
//
// Some open source application is distributed in the hope that it will
// be useful, but WITHOUT ANY WARRANTY; without even the implied warranty
// of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Foobar.  If not, see <http://www.gnu.org/licenses/>.
//
// @license GPL-3.0+ <http://spdx.org/licenses/GPL-3.0+>
//
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/miku/groupcover"
)

// Version displayed by application.
const Version = "0.0.10"

func main() {
	// TODO(miku): Adjust for AMSL format.
	// curl -s http://example.com/outboundservices/list?do=deduplication
	// [
	//   {
	//     "sourceID": "3",
	//     "sourceID_dedup": "0;108;120;16;17;18;19;26;4;59;63;72;74;84;86"
	//   },
	//   {
	//     "sourceID": "4",
	//     "sourceID_dedup": "0;108;120;16;17;18;19;26;63;72;74;84;86"
	//   },
	//   ...

	prefs := flag.String("prefs", "", "space separated string of preferences (most preferred first), e.g. 'B A C'")
	cpuprofile := flag.String("cpuprofile", "", "path to pprof output")
	verbose := flag.Bool("verbose", false, "more output")
	version := flag.Bool("version", false, "show version")
	column := flag.Int("f", 3, "column to use for grouping, one-based")
	lowerCase := flag.Bool("lower", false, "lowercase input")

	flag.Parse()

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *column < 1 {
		log.Fatal("column index must be non-negative")
	}

	groupcover.Verbose = *verbose

	var attrFunc groupcover.AttrFunc

	if *lowerCase {
		attrFunc = groupcover.ColumnLower(*column - 1)
	} else {
		attrFunc = groupcover.Column(*column - 1)
	}

	preferences := groupcover.Preferences{}

	// Parse preferences, if given.
	if *prefs != "" {
		fields := strings.Fields(*prefs)
		if len(fields) == 0 {
			log.Fatal("prefs must not be empty")
		}
		// Adjust the default ChoiceFunc.
		preferences.Default = groupcover.ListChooser(fields)
	}

	// A simple rewriter, that considers per-key preferences.
	rewriter := groupcover.SimpleRewriter(preferences)

	// Read from stdin, write to stdout, use third column as grouping criteria
	// and rewriter as rewriter.
	if err := groupcover.GroupRewrite(os.Stdin, os.Stdout, attrFunc, rewriter); err != nil {
		log.Fatal(err)
	}
}
