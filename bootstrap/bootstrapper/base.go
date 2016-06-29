// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package bootstrapper

import (
	"sync"

	"github.com/m3db/m3db/bootstrap"
	"github.com/m3db/m3db/interfaces/m3db"
	xtime "github.com/m3db/m3db/x/time"
)

const (
	baseBootstrapperName = "base"
)

// baseBootstrapper provides a skeleton for the interface methods.
type baseBootstrapper struct {
	s      m3db.Source
	dbOpts m3db.DatabaseOptions
	next   m3db.Bootstrapper
}

// NewBaseBootstrapper creates a new base bootstrapper.
func NewBaseBootstrapper(
	s m3db.Source,
	dbOpts m3db.DatabaseOptions,
	next m3db.Bootstrapper,
) m3db.Bootstrapper {
	bs := next
	if next == nil {
		bs = defaultNoOpBootstrapper
	}
	return &baseBootstrapper{s: s, dbOpts: dbOpts, next: bs}
}

// Bootstrap performs bootstrapping for the given shards and the associated time ranges.
func (bsb *baseBootstrapper) Bootstrap(shard uint32, targetRanges xtime.Ranges) (m3db.ShardResult, xtime.Ranges) {
	if xtime.IsEmpty(targetRanges) {
		return nil, nil
	}

	availableRanges := bsb.s.GetAvailability(shard, targetRanges)
	remainingRanges := targetRanges.RemoveRanges(availableRanges)

	var (
		wg                              sync.WaitGroup
		curResult, nextResult           m3db.ShardResult
		curUnfulfilled, nextUnfulfilled xtime.Ranges
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		nextResult, nextUnfulfilled = bsb.next.Bootstrap(shard, remainingRanges)
	}()

	curResult, curUnfulfilled = bsb.s.ReadData(shard, availableRanges)
	wg.Wait()

	mergedResults := bsb.mergeResults(curResult, nextResult)

	// If there are some time ranges the current bootstrapper can't fulfill,
	// pass it along to the next bootstrapper.
	if !xtime.IsEmpty(curUnfulfilled) {
		curResult, curUnfulfilled = bsb.next.Bootstrap(shard, curUnfulfilled)
		mergedResults = bsb.mergeResults(mergedResults, curResult)
	}

	mergedUnfulfilled := mergeTimeRanges(curUnfulfilled, nextUnfulfilled)
	return mergedResults, mergedUnfulfilled
}

func (bsb *baseBootstrapper) mergeResults(results ...m3db.ShardResult) m3db.ShardResult {
	final := bootstrap.NewShardResult(bsb.dbOpts)
	for _, result := range results {
		final.AddResult(result)
	}
	return final
}

func mergeTimeRanges(ranges ...xtime.Ranges) xtime.Ranges {
	final := xtime.NewRanges()
	for _, tr := range ranges {
		final = final.AddRanges(tr)
	}
	return final
}

// String returns the name of the bootstrapper.
func (bsb *baseBootstrapper) String() string {
	return baseBootstrapperName
}