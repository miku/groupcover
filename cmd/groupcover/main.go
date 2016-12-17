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
	"log"
	"os"

	"github.com/miku/groupcover"
)

func main() {
	// Preference per key.

	// preferences := groupcover.PreferenceMap{
	// 	"K1": groupcover.LexChoice,
	// 	"K2": groupcover.LexChoice,
	// 	"K3": groupcover.LexChoice,
	// }

	// preferences := groupcover.PreferenceMap{
	// 	"DE-105":    groupcover.LexChoice,
	// 	"DE-14":     groupcover.LexChoice,
	// 	"DE-15":     groupcover.LexChoice,
	// 	"DE-15-FID": groupcover.LexChoice,
	// 	"DE-1972":   groupcover.LexChoice,
	// 	"DE-540":    groupcover.LexChoice,
	// 	"DE-Bn3":    groupcover.LexChoice,
	// 	"DE-Brt1":   groupcover.LexChoice,
	// 	"DE-Ch1":    groupcover.LexChoice,
	// 	"DE-D13":    groupcover.LexChoice,
	// 	"DE-D161":   groupcover.LexChoice,
	// 	"DE-Gla1":   groupcover.LexChoice,
	// 	"DE-Kn38":   groupcover.LexChoice,
	// 	"DE-L152":   groupcover.LexChoice,
	// 	"DE-L242":   groupcover.LexChoice,
	// 	"DE-Zi4":    groupcover.LexChoice,
	// 	"DE-Zwi2":   groupcover.LexChoice,
	// }

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

	// Use the third column as grouping criteria.
	thirdColumn := groupcover.Column(2)

	// If not set explicitly, defaults to lexicographic order of the group identifiers.
	preferences := groupcover.PreferenceMap{}
	// A simple rewriter, that considers per-key preferences.
	rewriter := groupcover.SimpleRewriter(preferences)

	// Read from stdin, write to stdout, use third column as grouping criteria
	// and rewriter as rewriter.
	if err := groupcover.GroupRewrite(os.Stdin, os.Stdout, thirdColumn, rewriter); err != nil {
		log.Fatal(err)
	}
}
