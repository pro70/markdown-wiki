package version

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/irgangla/markdown-wiki/log"
	"github.com/irgangla/markdown-wiki/sdk"
)

var (
	commitName      string
	commitMail      string
	pushAfterCommit bool
)

// Start versioning
func Start(name, mail string, push bool) {
	log.Info("VERSION", "Versioning started.")
	sdk.UpdateEvents = make(chan sdk.Event, 5)
	commitName = name
	commitMail = mail
	pushAfterCommit = push
	go processEvents()
}

// Stop versioning
func Stop() {
	log.Info("VERSION", "Versioning stopped.")
	close(sdk.UpdateEvents)
}

func processEvents() {
	for e := range sdk.UpdateEvents {
		log.Info("VERSION", "Processing update event", e)

		r, err := git.PlainOpen(".")
		if err != nil {
			log.Error("VERSION", "open repository:", err)
			r, err = git.PlainInit(".", false)
			if err != nil {
				log.Error("VERSION", "create repository:", err)
				continue
			}
		}

		w, err := r.Worktree()
		if err != nil {
			log.Error("VERSION", "read worktree:", err)
			continue
		}

		message := fmt.Sprintf("update markdown file %v", e.Data)
		commit, err := w.Commit(message, &git.CommitOptions{
			Author: &object.Signature{
				Name:  commitName,
				Email: commitMail,
				When:  time.Now(),
			},
		})
		if err != nil {
			log.Error("VERSION", "commit change:", err)
			continue
		}

		obj, err := r.CommitObject(commit)
		if err != nil {
			log.Error("VERSION", "read commit:", err)
			continue
		}
		log.Info("VERSION", "update", obj)

		if pushAfterCommit {
			log.Info("VERSION", "push repository")
			err = r.Push(&git.PushOptions{})
			if err != nil {
				log.Error("VERSION", "push commit:", err)
			}
		} else {
			log.Debug("VERSION", "push after commit disabled")
		}
	}
}
