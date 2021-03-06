package main

import (
	"github.com/SpiderOak/gonimbusio/nimbusapi"

	"io"
	"os"
)

type WorkUnit struct {
	collection          string
	key                 string
	conjoinedIdentifier string
	conjoinedPart       int
	offset              int64
	size                int64
}

type WorkResult struct {
	workerId      int
	conjoinedPart int
	size          int64
	err           error
	action        string
}

func worker(id int, filePath string, requester nimbusapi.Requester,
	work <-chan WorkUnit, results chan<- WorkResult) {
	result := WorkResult{}
	result.workerId = id

	file, err := os.Open(filePath)
	if err != nil {
		result.err = err
		result.action = "Open"
		results <- result
		return
	}
	defer file.Close()

	for workUnit := range work {

		_, err = file.Seek(workUnit.offset, 0)
		if err != nil {
			result.err = err
			result.action = "Seek"
			results <- result
			return
		}

		conjoinedParams := nimbusapi.ConjoinedParams{
			workUnit.conjoinedIdentifier, workUnit.conjoinedPart}

		_, err := nimbusapi.Archive(requester, workUnit.collection,
			workUnit.key, &conjoinedParams,
			io.LimitReader(file, workUnit.size))

		if err != nil {
			result.err = err
			result.action = "Archive"
			results <- result
			return
		}

		result.conjoinedPart = workUnit.conjoinedPart
		result.size = workUnit.size
		results <- result
	}
}
