package clarum

import (
	"fmt"
	"github.com/goclarum/clarum/core/control"
	"github.com/goclarum/clarum/http"
)

func Http() http.Builder {
	return &http.EndpointBuilder{}
}

func Finish() {
	fmt.Println(fmt.Sprintf("Waiting for all actions to finish."))

	control.RunningActions.Wait()

	fmt.Println(fmt.Sprintf("All actions finished."))
}
