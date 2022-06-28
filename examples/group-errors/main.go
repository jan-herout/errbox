/**************
 * Output:
 **************

Got 2 errors:
----------------------------
# 1
bad stuff happened
 +--> because we were careless
    @ /main.go:29 (main)

----------------------------
# 2
after that, another bad thing happened
 +--> karma!
    @ /main.go:30 (main)
*/
package main

import (
	"fmt"

	"github.com/jan-herout/errbox"
)

func main() {
	errbox.OmitPrefixFromTrace("C:/Git/errbox/examples/group-errors") // Cleanup trace, do not show full paths.
	var be errbox.Box
	be.PushIf(fmt.Errorf("bad stuff happened"), "because we were careless")
	if be.PushIf(fmt.Errorf("after that, another bad thing happened"), "karma!") {
		fmt.Println(be)
	}
}
