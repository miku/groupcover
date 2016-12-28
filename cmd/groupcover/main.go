//  Copyright 2016 by Leipzig University Library, http://ub.uni-leipzig.de
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

const Version = "0.0.5"

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

	groupcover.Verbose = *verbose

	// Use the third column as grouping criteria.
	thirdColumn := groupcover.Column(2)

	// Parse preferences.
	if *prefs != "" {
		fields := strings.Fields(*prefs)
		if len(fields) == 0 {
			log.Fatal("preference must not be empty")
		}
		// It would be better to not rely on package vars, but ok for now.
		groupcover.DefaultChoiceFunc = groupcover.ListChooser(fields)
	}

	// A simple rewriter, that considers per-key preferences. First, test with
	// the same default for all keys (groupcover.DefaultChoiceFunc).
	rewriter := groupcover.SimpleRewriter(groupcover.PreferenceMap{})

	// Read from stdin, write to stdout, use third column as grouping criteria
	// and rewriter as rewriter.
	if err := groupcover.GroupRewrite(os.Stdin, os.Stdout, thirdColumn, rewriter); err != nil {
		log.Fatal(err)
	}
}
