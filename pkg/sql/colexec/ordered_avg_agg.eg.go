// Code generated by execgen; DO NOT EDIT.
// Copyright 2018 The Cockroach Authors.
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package colexec

import (
	"unsafe"

	"github.com/cockroachdb/apd/v2"
	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/sql/colexecbase/colexecerror"
	"github.com/cockroachdb/cockroach/pkg/sql/colmem"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/types"
	"github.com/cockroachdb/cockroach/pkg/util/duration"
	"github.com/cockroachdb/errors"
)

func newAvgOrderedAggAlloc(
	allocator *colmem.Allocator, t *types.T, allocSize int64,
) (aggregateFuncAlloc, error) {
	allocBase := aggAllocBase{allocator: allocator, allocSize: allocSize}
	switch t.Family() {
	case types.IntFamily:
		switch t.Width() {
		case 16:
			return &avgInt16OrderedAggAlloc{aggAllocBase: allocBase}, nil
		case 32:
			return &avgInt32OrderedAggAlloc{aggAllocBase: allocBase}, nil
		default:
			return &avgInt64OrderedAggAlloc{aggAllocBase: allocBase}, nil
		}
	case types.DecimalFamily:
		return &avgDecimalOrderedAggAlloc{aggAllocBase: allocBase}, nil
	case types.FloatFamily:
		return &avgFloat64OrderedAggAlloc{aggAllocBase: allocBase}, nil
	case types.IntervalFamily:
		return &avgIntervalOrderedAggAlloc{aggAllocBase: allocBase}, nil
	default:
		return nil, errors.Errorf("unsupported avg agg type %s", t.Name())
	}
}

type avgInt16OrderedAgg struct {
	orderedAggregateFuncBase
	scratch struct {
		// curSum keeps track of the sum of elements belonging to the current group,
		// so we can index into the slice once per group, instead of on each
		// iteration.
		curSum apd.Decimal
		// curCount keeps track of the number of elements that we've seen
		// belonging to the current group.
		curCount int64
		// vec points to the output vector.
		vec []apd.Decimal
		// foundNonNullForCurrentGroup tracks if we have seen any non-null values
		// for the group that is currently being aggregated.
		foundNonNullForCurrentGroup bool
	}
	overloadHelper overloadHelper
}

var _ aggregateFunc = &avgInt16OrderedAgg{}

func (a *avgInt16OrderedAgg) Init(groups []bool, vec coldata.Vec) {
	a.orderedAggregateFuncBase.Init(groups, vec)
	a.scratch.vec = vec.Decimal()
	a.Reset()
}

func (a *avgInt16OrderedAgg) Reset() {
	a.orderedAggregateFuncBase.Reset()
	a.scratch.curSum = zeroDecimalValue
	a.scratch.curCount = 0
	a.scratch.foundNonNullForCurrentGroup = false
}

func (a *avgInt16OrderedAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	// In order to inline the templated code of overloads, we need to have a
	// "_overloadHelper" local variable of type "overloadHelper".
	_overloadHelper := a.overloadHelper
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Int16(), vec.Nulls()
	groups := a.groups
	if sel == nil {
		_ = groups[inputLen-1]
		col = col[:inputLen]
		if nulls.MaybeHasNulls() {
			for i := range col {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {

						a.scratch.vec[a.curIdx].SetInt64(a.scratch.curCount)
						if _, err := tree.DecimalCtx.Quo(&a.scratch.vec[a.curIdx], &a.scratch.curSum, &a.scratch.vec[a.curIdx]); err != nil {
							colexecerror.InternalError(err)
						}
					}
					a.curIdx++
					a.scratch.curSum = zeroDecimalValue
					a.scratch.curCount = 0

					a.scratch.foundNonNullForCurrentGroup = false
				}

				var isNull bool
				isNull = nulls.NullAt(i)
				if !isNull {

					{

						tmpDec := &_overloadHelper.tmpDec1
						tmpDec.SetInt64(int64(col[i]))
						if _, err := tree.ExactCtx.Add(&a.scratch.curSum, &a.scratch.curSum, tmpDec); err != nil {
							colexecerror.ExpectedError(err)
						}
					}

					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		} else {
			for i := range col {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {

						a.scratch.vec[a.curIdx].SetInt64(a.scratch.curCount)
						if _, err := tree.DecimalCtx.Quo(&a.scratch.vec[a.curIdx], &a.scratch.curSum, &a.scratch.vec[a.curIdx]); err != nil {
							colexecerror.InternalError(err)
						}
					}
					a.curIdx++
					a.scratch.curSum = zeroDecimalValue
					a.scratch.curCount = 0

				}

				var isNull bool
				isNull = false
				if !isNull {

					{

						tmpDec := &_overloadHelper.tmpDec1
						tmpDec.SetInt64(int64(col[i]))
						if _, err := tree.ExactCtx.Add(&a.scratch.curSum, &a.scratch.curSum, tmpDec); err != nil {
							colexecerror.ExpectedError(err)
						}
					}

					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		}
	} else {
		sel = sel[:inputLen]
		if nulls.MaybeHasNulls() {
			for _, i := range sel {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {

						a.scratch.vec[a.curIdx].SetInt64(a.scratch.curCount)
						if _, err := tree.DecimalCtx.Quo(&a.scratch.vec[a.curIdx], &a.scratch.curSum, &a.scratch.vec[a.curIdx]); err != nil {
							colexecerror.InternalError(err)
						}
					}
					a.curIdx++
					a.scratch.curSum = zeroDecimalValue
					a.scratch.curCount = 0

					a.scratch.foundNonNullForCurrentGroup = false
				}

				var isNull bool
				isNull = nulls.NullAt(i)
				if !isNull {

					{

						tmpDec := &_overloadHelper.tmpDec1
						tmpDec.SetInt64(int64(col[i]))
						if _, err := tree.ExactCtx.Add(&a.scratch.curSum, &a.scratch.curSum, tmpDec); err != nil {
							colexecerror.ExpectedError(err)
						}
					}

					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		} else {
			for _, i := range sel {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {

						a.scratch.vec[a.curIdx].SetInt64(a.scratch.curCount)
						if _, err := tree.DecimalCtx.Quo(&a.scratch.vec[a.curIdx], &a.scratch.curSum, &a.scratch.vec[a.curIdx]); err != nil {
							colexecerror.InternalError(err)
						}
					}
					a.curIdx++
					a.scratch.curSum = zeroDecimalValue
					a.scratch.curCount = 0

				}

				var isNull bool
				isNull = false
				if !isNull {

					{

						tmpDec := &_overloadHelper.tmpDec1
						tmpDec.SetInt64(int64(col[i]))
						if _, err := tree.ExactCtx.Add(&a.scratch.curSum, &a.scratch.curSum, tmpDec); err != nil {
							colexecerror.ExpectedError(err)
						}
					}

					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		}
	}
}

func (a *avgInt16OrderedAgg) Flush(outputIdx int) {
	// The aggregation is finished. Flush the last value. If we haven't found
	// any non-nulls for this group so far, the output for this group should be
	// NULL.
	// Go around "argument overwritten before first use" linter error.
	_ = outputIdx
	outputIdx = a.curIdx
	a.curIdx++
	if !a.scratch.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {

		a.scratch.vec[outputIdx].SetInt64(a.scratch.curCount)
		if _, err := tree.DecimalCtx.Quo(&a.scratch.vec[outputIdx], &a.scratch.curSum, &a.scratch.vec[outputIdx]); err != nil {
			colexecerror.InternalError(err)
		}
	}
}

type avgInt16OrderedAggAlloc struct {
	aggAllocBase
	aggFuncs []avgInt16OrderedAgg
}

var _ aggregateFuncAlloc = &avgInt16OrderedAggAlloc{}

const sizeOfAvgInt16OrderedAgg = int64(unsafe.Sizeof(avgInt16OrderedAgg{}))
const avgInt16OrderedAggSliceOverhead = int64(unsafe.Sizeof([]avgInt16OrderedAgg{}))

func (a *avgInt16OrderedAggAlloc) newAggFunc() aggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(avgInt16OrderedAggSliceOverhead + sizeOfAvgInt16OrderedAgg*a.allocSize)
		a.aggFuncs = make([]avgInt16OrderedAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	a.aggFuncs = a.aggFuncs[1:]
	return f
}

type avgInt32OrderedAgg struct {
	orderedAggregateFuncBase
	scratch struct {
		// curSum keeps track of the sum of elements belonging to the current group,
		// so we can index into the slice once per group, instead of on each
		// iteration.
		curSum apd.Decimal
		// curCount keeps track of the number of elements that we've seen
		// belonging to the current group.
		curCount int64
		// vec points to the output vector.
		vec []apd.Decimal
		// foundNonNullForCurrentGroup tracks if we have seen any non-null values
		// for the group that is currently being aggregated.
		foundNonNullForCurrentGroup bool
	}
	overloadHelper overloadHelper
}

var _ aggregateFunc = &avgInt32OrderedAgg{}

func (a *avgInt32OrderedAgg) Init(groups []bool, vec coldata.Vec) {
	a.orderedAggregateFuncBase.Init(groups, vec)
	a.scratch.vec = vec.Decimal()
	a.Reset()
}

func (a *avgInt32OrderedAgg) Reset() {
	a.orderedAggregateFuncBase.Reset()
	a.scratch.curSum = zeroDecimalValue
	a.scratch.curCount = 0
	a.scratch.foundNonNullForCurrentGroup = false
}

func (a *avgInt32OrderedAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	// In order to inline the templated code of overloads, we need to have a
	// "_overloadHelper" local variable of type "overloadHelper".
	_overloadHelper := a.overloadHelper
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Int32(), vec.Nulls()
	groups := a.groups
	if sel == nil {
		_ = groups[inputLen-1]
		col = col[:inputLen]
		if nulls.MaybeHasNulls() {
			for i := range col {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {

						a.scratch.vec[a.curIdx].SetInt64(a.scratch.curCount)
						if _, err := tree.DecimalCtx.Quo(&a.scratch.vec[a.curIdx], &a.scratch.curSum, &a.scratch.vec[a.curIdx]); err != nil {
							colexecerror.InternalError(err)
						}
					}
					a.curIdx++
					a.scratch.curSum = zeroDecimalValue
					a.scratch.curCount = 0

					a.scratch.foundNonNullForCurrentGroup = false
				}

				var isNull bool
				isNull = nulls.NullAt(i)
				if !isNull {

					{

						tmpDec := &_overloadHelper.tmpDec1
						tmpDec.SetInt64(int64(col[i]))
						if _, err := tree.ExactCtx.Add(&a.scratch.curSum, &a.scratch.curSum, tmpDec); err != nil {
							colexecerror.ExpectedError(err)
						}
					}

					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		} else {
			for i := range col {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {

						a.scratch.vec[a.curIdx].SetInt64(a.scratch.curCount)
						if _, err := tree.DecimalCtx.Quo(&a.scratch.vec[a.curIdx], &a.scratch.curSum, &a.scratch.vec[a.curIdx]); err != nil {
							colexecerror.InternalError(err)
						}
					}
					a.curIdx++
					a.scratch.curSum = zeroDecimalValue
					a.scratch.curCount = 0

				}

				var isNull bool
				isNull = false
				if !isNull {

					{

						tmpDec := &_overloadHelper.tmpDec1
						tmpDec.SetInt64(int64(col[i]))
						if _, err := tree.ExactCtx.Add(&a.scratch.curSum, &a.scratch.curSum, tmpDec); err != nil {
							colexecerror.ExpectedError(err)
						}
					}

					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		}
	} else {
		sel = sel[:inputLen]
		if nulls.MaybeHasNulls() {
			for _, i := range sel {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {

						a.scratch.vec[a.curIdx].SetInt64(a.scratch.curCount)
						if _, err := tree.DecimalCtx.Quo(&a.scratch.vec[a.curIdx], &a.scratch.curSum, &a.scratch.vec[a.curIdx]); err != nil {
							colexecerror.InternalError(err)
						}
					}
					a.curIdx++
					a.scratch.curSum = zeroDecimalValue
					a.scratch.curCount = 0

					a.scratch.foundNonNullForCurrentGroup = false
				}

				var isNull bool
				isNull = nulls.NullAt(i)
				if !isNull {

					{

						tmpDec := &_overloadHelper.tmpDec1
						tmpDec.SetInt64(int64(col[i]))
						if _, err := tree.ExactCtx.Add(&a.scratch.curSum, &a.scratch.curSum, tmpDec); err != nil {
							colexecerror.ExpectedError(err)
						}
					}

					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		} else {
			for _, i := range sel {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {

						a.scratch.vec[a.curIdx].SetInt64(a.scratch.curCount)
						if _, err := tree.DecimalCtx.Quo(&a.scratch.vec[a.curIdx], &a.scratch.curSum, &a.scratch.vec[a.curIdx]); err != nil {
							colexecerror.InternalError(err)
						}
					}
					a.curIdx++
					a.scratch.curSum = zeroDecimalValue
					a.scratch.curCount = 0

				}

				var isNull bool
				isNull = false
				if !isNull {

					{

						tmpDec := &_overloadHelper.tmpDec1
						tmpDec.SetInt64(int64(col[i]))
						if _, err := tree.ExactCtx.Add(&a.scratch.curSum, &a.scratch.curSum, tmpDec); err != nil {
							colexecerror.ExpectedError(err)
						}
					}

					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		}
	}
}

func (a *avgInt32OrderedAgg) Flush(outputIdx int) {
	// The aggregation is finished. Flush the last value. If we haven't found
	// any non-nulls for this group so far, the output for this group should be
	// NULL.
	// Go around "argument overwritten before first use" linter error.
	_ = outputIdx
	outputIdx = a.curIdx
	a.curIdx++
	if !a.scratch.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {

		a.scratch.vec[outputIdx].SetInt64(a.scratch.curCount)
		if _, err := tree.DecimalCtx.Quo(&a.scratch.vec[outputIdx], &a.scratch.curSum, &a.scratch.vec[outputIdx]); err != nil {
			colexecerror.InternalError(err)
		}
	}
}

type avgInt32OrderedAggAlloc struct {
	aggAllocBase
	aggFuncs []avgInt32OrderedAgg
}

var _ aggregateFuncAlloc = &avgInt32OrderedAggAlloc{}

const sizeOfAvgInt32OrderedAgg = int64(unsafe.Sizeof(avgInt32OrderedAgg{}))
const avgInt32OrderedAggSliceOverhead = int64(unsafe.Sizeof([]avgInt32OrderedAgg{}))

func (a *avgInt32OrderedAggAlloc) newAggFunc() aggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(avgInt32OrderedAggSliceOverhead + sizeOfAvgInt32OrderedAgg*a.allocSize)
		a.aggFuncs = make([]avgInt32OrderedAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	a.aggFuncs = a.aggFuncs[1:]
	return f
}

type avgInt64OrderedAgg struct {
	orderedAggregateFuncBase
	scratch struct {
		// curSum keeps track of the sum of elements belonging to the current group,
		// so we can index into the slice once per group, instead of on each
		// iteration.
		curSum apd.Decimal
		// curCount keeps track of the number of elements that we've seen
		// belonging to the current group.
		curCount int64
		// vec points to the output vector.
		vec []apd.Decimal
		// foundNonNullForCurrentGroup tracks if we have seen any non-null values
		// for the group that is currently being aggregated.
		foundNonNullForCurrentGroup bool
	}
	overloadHelper overloadHelper
}

var _ aggregateFunc = &avgInt64OrderedAgg{}

func (a *avgInt64OrderedAgg) Init(groups []bool, vec coldata.Vec) {
	a.orderedAggregateFuncBase.Init(groups, vec)
	a.scratch.vec = vec.Decimal()
	a.Reset()
}

func (a *avgInt64OrderedAgg) Reset() {
	a.orderedAggregateFuncBase.Reset()
	a.scratch.curSum = zeroDecimalValue
	a.scratch.curCount = 0
	a.scratch.foundNonNullForCurrentGroup = false
}

func (a *avgInt64OrderedAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	// In order to inline the templated code of overloads, we need to have a
	// "_overloadHelper" local variable of type "overloadHelper".
	_overloadHelper := a.overloadHelper
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Int64(), vec.Nulls()
	groups := a.groups
	if sel == nil {
		_ = groups[inputLen-1]
		col = col[:inputLen]
		if nulls.MaybeHasNulls() {
			for i := range col {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {

						a.scratch.vec[a.curIdx].SetInt64(a.scratch.curCount)
						if _, err := tree.DecimalCtx.Quo(&a.scratch.vec[a.curIdx], &a.scratch.curSum, &a.scratch.vec[a.curIdx]); err != nil {
							colexecerror.InternalError(err)
						}
					}
					a.curIdx++
					a.scratch.curSum = zeroDecimalValue
					a.scratch.curCount = 0

					a.scratch.foundNonNullForCurrentGroup = false
				}

				var isNull bool
				isNull = nulls.NullAt(i)
				if !isNull {

					{

						tmpDec := &_overloadHelper.tmpDec1
						tmpDec.SetInt64(int64(col[i]))
						if _, err := tree.ExactCtx.Add(&a.scratch.curSum, &a.scratch.curSum, tmpDec); err != nil {
							colexecerror.ExpectedError(err)
						}
					}

					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		} else {
			for i := range col {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {

						a.scratch.vec[a.curIdx].SetInt64(a.scratch.curCount)
						if _, err := tree.DecimalCtx.Quo(&a.scratch.vec[a.curIdx], &a.scratch.curSum, &a.scratch.vec[a.curIdx]); err != nil {
							colexecerror.InternalError(err)
						}
					}
					a.curIdx++
					a.scratch.curSum = zeroDecimalValue
					a.scratch.curCount = 0

				}

				var isNull bool
				isNull = false
				if !isNull {

					{

						tmpDec := &_overloadHelper.tmpDec1
						tmpDec.SetInt64(int64(col[i]))
						if _, err := tree.ExactCtx.Add(&a.scratch.curSum, &a.scratch.curSum, tmpDec); err != nil {
							colexecerror.ExpectedError(err)
						}
					}

					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		}
	} else {
		sel = sel[:inputLen]
		if nulls.MaybeHasNulls() {
			for _, i := range sel {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {

						a.scratch.vec[a.curIdx].SetInt64(a.scratch.curCount)
						if _, err := tree.DecimalCtx.Quo(&a.scratch.vec[a.curIdx], &a.scratch.curSum, &a.scratch.vec[a.curIdx]); err != nil {
							colexecerror.InternalError(err)
						}
					}
					a.curIdx++
					a.scratch.curSum = zeroDecimalValue
					a.scratch.curCount = 0

					a.scratch.foundNonNullForCurrentGroup = false
				}

				var isNull bool
				isNull = nulls.NullAt(i)
				if !isNull {

					{

						tmpDec := &_overloadHelper.tmpDec1
						tmpDec.SetInt64(int64(col[i]))
						if _, err := tree.ExactCtx.Add(&a.scratch.curSum, &a.scratch.curSum, tmpDec); err != nil {
							colexecerror.ExpectedError(err)
						}
					}

					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		} else {
			for _, i := range sel {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {

						a.scratch.vec[a.curIdx].SetInt64(a.scratch.curCount)
						if _, err := tree.DecimalCtx.Quo(&a.scratch.vec[a.curIdx], &a.scratch.curSum, &a.scratch.vec[a.curIdx]); err != nil {
							colexecerror.InternalError(err)
						}
					}
					a.curIdx++
					a.scratch.curSum = zeroDecimalValue
					a.scratch.curCount = 0

				}

				var isNull bool
				isNull = false
				if !isNull {

					{

						tmpDec := &_overloadHelper.tmpDec1
						tmpDec.SetInt64(int64(col[i]))
						if _, err := tree.ExactCtx.Add(&a.scratch.curSum, &a.scratch.curSum, tmpDec); err != nil {
							colexecerror.ExpectedError(err)
						}
					}

					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		}
	}
}

func (a *avgInt64OrderedAgg) Flush(outputIdx int) {
	// The aggregation is finished. Flush the last value. If we haven't found
	// any non-nulls for this group so far, the output for this group should be
	// NULL.
	// Go around "argument overwritten before first use" linter error.
	_ = outputIdx
	outputIdx = a.curIdx
	a.curIdx++
	if !a.scratch.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {

		a.scratch.vec[outputIdx].SetInt64(a.scratch.curCount)
		if _, err := tree.DecimalCtx.Quo(&a.scratch.vec[outputIdx], &a.scratch.curSum, &a.scratch.vec[outputIdx]); err != nil {
			colexecerror.InternalError(err)
		}
	}
}

type avgInt64OrderedAggAlloc struct {
	aggAllocBase
	aggFuncs []avgInt64OrderedAgg
}

var _ aggregateFuncAlloc = &avgInt64OrderedAggAlloc{}

const sizeOfAvgInt64OrderedAgg = int64(unsafe.Sizeof(avgInt64OrderedAgg{}))
const avgInt64OrderedAggSliceOverhead = int64(unsafe.Sizeof([]avgInt64OrderedAgg{}))

func (a *avgInt64OrderedAggAlloc) newAggFunc() aggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(avgInt64OrderedAggSliceOverhead + sizeOfAvgInt64OrderedAgg*a.allocSize)
		a.aggFuncs = make([]avgInt64OrderedAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	a.aggFuncs = a.aggFuncs[1:]
	return f
}

type avgDecimalOrderedAgg struct {
	orderedAggregateFuncBase
	scratch struct {
		// curSum keeps track of the sum of elements belonging to the current group,
		// so we can index into the slice once per group, instead of on each
		// iteration.
		curSum apd.Decimal
		// curCount keeps track of the number of elements that we've seen
		// belonging to the current group.
		curCount int64
		// vec points to the output vector.
		vec []apd.Decimal
		// foundNonNullForCurrentGroup tracks if we have seen any non-null values
		// for the group that is currently being aggregated.
		foundNonNullForCurrentGroup bool
	}
}

var _ aggregateFunc = &avgDecimalOrderedAgg{}

func (a *avgDecimalOrderedAgg) Init(groups []bool, vec coldata.Vec) {
	a.orderedAggregateFuncBase.Init(groups, vec)
	a.scratch.vec = vec.Decimal()
	a.Reset()
}

func (a *avgDecimalOrderedAgg) Reset() {
	a.orderedAggregateFuncBase.Reset()
	a.scratch.curSum = zeroDecimalValue
	a.scratch.curCount = 0
	a.scratch.foundNonNullForCurrentGroup = false
}

func (a *avgDecimalOrderedAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Decimal(), vec.Nulls()
	groups := a.groups
	if sel == nil {
		_ = groups[inputLen-1]
		col = col[:inputLen]
		if nulls.MaybeHasNulls() {
			for i := range col {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {

						a.scratch.vec[a.curIdx].SetInt64(a.scratch.curCount)
						if _, err := tree.DecimalCtx.Quo(&a.scratch.vec[a.curIdx], &a.scratch.curSum, &a.scratch.vec[a.curIdx]); err != nil {
							colexecerror.InternalError(err)
						}
					}
					a.curIdx++
					a.scratch.curSum = zeroDecimalValue
					a.scratch.curCount = 0

					a.scratch.foundNonNullForCurrentGroup = false
				}

				var isNull bool
				isNull = nulls.NullAt(i)
				if !isNull {

					{

						_, err := tree.ExactCtx.Add(&a.scratch.curSum, &a.scratch.curSum, &col[i])
						if err != nil {
							colexecerror.ExpectedError(err)
						}
					}

					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		} else {
			for i := range col {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {

						a.scratch.vec[a.curIdx].SetInt64(a.scratch.curCount)
						if _, err := tree.DecimalCtx.Quo(&a.scratch.vec[a.curIdx], &a.scratch.curSum, &a.scratch.vec[a.curIdx]); err != nil {
							colexecerror.InternalError(err)
						}
					}
					a.curIdx++
					a.scratch.curSum = zeroDecimalValue
					a.scratch.curCount = 0

				}

				var isNull bool
				isNull = false
				if !isNull {

					{

						_, err := tree.ExactCtx.Add(&a.scratch.curSum, &a.scratch.curSum, &col[i])
						if err != nil {
							colexecerror.ExpectedError(err)
						}
					}

					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		}
	} else {
		sel = sel[:inputLen]
		if nulls.MaybeHasNulls() {
			for _, i := range sel {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {

						a.scratch.vec[a.curIdx].SetInt64(a.scratch.curCount)
						if _, err := tree.DecimalCtx.Quo(&a.scratch.vec[a.curIdx], &a.scratch.curSum, &a.scratch.vec[a.curIdx]); err != nil {
							colexecerror.InternalError(err)
						}
					}
					a.curIdx++
					a.scratch.curSum = zeroDecimalValue
					a.scratch.curCount = 0

					a.scratch.foundNonNullForCurrentGroup = false
				}

				var isNull bool
				isNull = nulls.NullAt(i)
				if !isNull {

					{

						_, err := tree.ExactCtx.Add(&a.scratch.curSum, &a.scratch.curSum, &col[i])
						if err != nil {
							colexecerror.ExpectedError(err)
						}
					}

					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		} else {
			for _, i := range sel {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {

						a.scratch.vec[a.curIdx].SetInt64(a.scratch.curCount)
						if _, err := tree.DecimalCtx.Quo(&a.scratch.vec[a.curIdx], &a.scratch.curSum, &a.scratch.vec[a.curIdx]); err != nil {
							colexecerror.InternalError(err)
						}
					}
					a.curIdx++
					a.scratch.curSum = zeroDecimalValue
					a.scratch.curCount = 0

				}

				var isNull bool
				isNull = false
				if !isNull {

					{

						_, err := tree.ExactCtx.Add(&a.scratch.curSum, &a.scratch.curSum, &col[i])
						if err != nil {
							colexecerror.ExpectedError(err)
						}
					}

					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		}
	}
}

func (a *avgDecimalOrderedAgg) Flush(outputIdx int) {
	// The aggregation is finished. Flush the last value. If we haven't found
	// any non-nulls for this group so far, the output for this group should be
	// NULL.
	// Go around "argument overwritten before first use" linter error.
	_ = outputIdx
	outputIdx = a.curIdx
	a.curIdx++
	if !a.scratch.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {

		a.scratch.vec[outputIdx].SetInt64(a.scratch.curCount)
		if _, err := tree.DecimalCtx.Quo(&a.scratch.vec[outputIdx], &a.scratch.curSum, &a.scratch.vec[outputIdx]); err != nil {
			colexecerror.InternalError(err)
		}
	}
}

type avgDecimalOrderedAggAlloc struct {
	aggAllocBase
	aggFuncs []avgDecimalOrderedAgg
}

var _ aggregateFuncAlloc = &avgDecimalOrderedAggAlloc{}

const sizeOfAvgDecimalOrderedAgg = int64(unsafe.Sizeof(avgDecimalOrderedAgg{}))
const avgDecimalOrderedAggSliceOverhead = int64(unsafe.Sizeof([]avgDecimalOrderedAgg{}))

func (a *avgDecimalOrderedAggAlloc) newAggFunc() aggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(avgDecimalOrderedAggSliceOverhead + sizeOfAvgDecimalOrderedAgg*a.allocSize)
		a.aggFuncs = make([]avgDecimalOrderedAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	a.aggFuncs = a.aggFuncs[1:]
	return f
}

type avgFloat64OrderedAgg struct {
	orderedAggregateFuncBase
	scratch struct {
		// curSum keeps track of the sum of elements belonging to the current group,
		// so we can index into the slice once per group, instead of on each
		// iteration.
		curSum float64
		// curCount keeps track of the number of elements that we've seen
		// belonging to the current group.
		curCount int64
		// vec points to the output vector.
		vec []float64
		// foundNonNullForCurrentGroup tracks if we have seen any non-null values
		// for the group that is currently being aggregated.
		foundNonNullForCurrentGroup bool
	}
}

var _ aggregateFunc = &avgFloat64OrderedAgg{}

func (a *avgFloat64OrderedAgg) Init(groups []bool, vec coldata.Vec) {
	a.orderedAggregateFuncBase.Init(groups, vec)
	a.scratch.vec = vec.Float64()
	a.Reset()
}

func (a *avgFloat64OrderedAgg) Reset() {
	a.orderedAggregateFuncBase.Reset()
	a.scratch.curSum = zeroFloat64Value
	a.scratch.curCount = 0
	a.scratch.foundNonNullForCurrentGroup = false
}

func (a *avgFloat64OrderedAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Float64(), vec.Nulls()
	groups := a.groups
	if sel == nil {
		_ = groups[inputLen-1]
		col = col[:inputLen]
		if nulls.MaybeHasNulls() {
			for i := range col {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {
						a.scratch.vec[a.curIdx] = a.scratch.curSum / float64(a.scratch.curCount)
					}
					a.curIdx++
					a.scratch.curSum = zeroFloat64Value
					a.scratch.curCount = 0

					a.scratch.foundNonNullForCurrentGroup = false
				}

				var isNull bool
				isNull = nulls.NullAt(i)
				if !isNull {

					{

						a.scratch.curSum = float64(a.scratch.curSum) + float64(col[i])
					}

					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		} else {
			for i := range col {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {
						a.scratch.vec[a.curIdx] = a.scratch.curSum / float64(a.scratch.curCount)
					}
					a.curIdx++
					a.scratch.curSum = zeroFloat64Value
					a.scratch.curCount = 0

				}

				var isNull bool
				isNull = false
				if !isNull {

					{

						a.scratch.curSum = float64(a.scratch.curSum) + float64(col[i])
					}

					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		}
	} else {
		sel = sel[:inputLen]
		if nulls.MaybeHasNulls() {
			for _, i := range sel {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {
						a.scratch.vec[a.curIdx] = a.scratch.curSum / float64(a.scratch.curCount)
					}
					a.curIdx++
					a.scratch.curSum = zeroFloat64Value
					a.scratch.curCount = 0

					a.scratch.foundNonNullForCurrentGroup = false
				}

				var isNull bool
				isNull = nulls.NullAt(i)
				if !isNull {

					{

						a.scratch.curSum = float64(a.scratch.curSum) + float64(col[i])
					}

					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		} else {
			for _, i := range sel {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {
						a.scratch.vec[a.curIdx] = a.scratch.curSum / float64(a.scratch.curCount)
					}
					a.curIdx++
					a.scratch.curSum = zeroFloat64Value
					a.scratch.curCount = 0

				}

				var isNull bool
				isNull = false
				if !isNull {

					{

						a.scratch.curSum = float64(a.scratch.curSum) + float64(col[i])
					}

					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		}
	}
}

func (a *avgFloat64OrderedAgg) Flush(outputIdx int) {
	// The aggregation is finished. Flush the last value. If we haven't found
	// any non-nulls for this group so far, the output for this group should be
	// NULL.
	// Go around "argument overwritten before first use" linter error.
	_ = outputIdx
	outputIdx = a.curIdx
	a.curIdx++
	if !a.scratch.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {
		a.scratch.vec[outputIdx] = a.scratch.curSum / float64(a.scratch.curCount)
	}
}

type avgFloat64OrderedAggAlloc struct {
	aggAllocBase
	aggFuncs []avgFloat64OrderedAgg
}

var _ aggregateFuncAlloc = &avgFloat64OrderedAggAlloc{}

const sizeOfAvgFloat64OrderedAgg = int64(unsafe.Sizeof(avgFloat64OrderedAgg{}))
const avgFloat64OrderedAggSliceOverhead = int64(unsafe.Sizeof([]avgFloat64OrderedAgg{}))

func (a *avgFloat64OrderedAggAlloc) newAggFunc() aggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(avgFloat64OrderedAggSliceOverhead + sizeOfAvgFloat64OrderedAgg*a.allocSize)
		a.aggFuncs = make([]avgFloat64OrderedAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	a.aggFuncs = a.aggFuncs[1:]
	return f
}

type avgIntervalOrderedAgg struct {
	orderedAggregateFuncBase
	scratch struct {
		// curSum keeps track of the sum of elements belonging to the current group,
		// so we can index into the slice once per group, instead of on each
		// iteration.
		curSum duration.Duration
		// curCount keeps track of the number of elements that we've seen
		// belonging to the current group.
		curCount int64
		// vec points to the output vector.
		vec []duration.Duration
		// foundNonNullForCurrentGroup tracks if we have seen any non-null values
		// for the group that is currently being aggregated.
		foundNonNullForCurrentGroup bool
	}
}

var _ aggregateFunc = &avgIntervalOrderedAgg{}

func (a *avgIntervalOrderedAgg) Init(groups []bool, vec coldata.Vec) {
	a.orderedAggregateFuncBase.Init(groups, vec)
	a.scratch.vec = vec.Interval()
	a.Reset()
}

func (a *avgIntervalOrderedAgg) Reset() {
	a.orderedAggregateFuncBase.Reset()
	a.scratch.curSum = zeroIntervalValue
	a.scratch.curCount = 0
	a.scratch.foundNonNullForCurrentGroup = false
}

func (a *avgIntervalOrderedAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, inputLen int, sel []int,
) {
	vec := vecs[inputIdxs[0]]
	col, nulls := vec.Interval(), vec.Nulls()
	groups := a.groups
	if sel == nil {
		_ = groups[inputLen-1]
		col = col[:inputLen]
		if nulls.MaybeHasNulls() {
			for i := range col {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {
						a.scratch.vec[a.curIdx] = a.scratch.curSum.Div(int64(a.scratch.curCount))
					}
					a.curIdx++
					a.scratch.curSum = zeroIntervalValue
					a.scratch.curCount = 0

					a.scratch.foundNonNullForCurrentGroup = false
				}

				var isNull bool
				isNull = nulls.NullAt(i)
				if !isNull {
					a.scratch.curSum = a.scratch.curSum.Add(col[i])
					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		} else {
			for i := range col {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {
						a.scratch.vec[a.curIdx] = a.scratch.curSum.Div(int64(a.scratch.curCount))
					}
					a.curIdx++
					a.scratch.curSum = zeroIntervalValue
					a.scratch.curCount = 0

				}

				var isNull bool
				isNull = false
				if !isNull {
					a.scratch.curSum = a.scratch.curSum.Add(col[i])
					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		}
	} else {
		sel = sel[:inputLen]
		if nulls.MaybeHasNulls() {
			for _, i := range sel {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {
						a.scratch.vec[a.curIdx] = a.scratch.curSum.Div(int64(a.scratch.curCount))
					}
					a.curIdx++
					a.scratch.curSum = zeroIntervalValue
					a.scratch.curCount = 0

					a.scratch.foundNonNullForCurrentGroup = false
				}

				var isNull bool
				isNull = nulls.NullAt(i)
				if !isNull {
					a.scratch.curSum = a.scratch.curSum.Add(col[i])
					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		} else {
			for _, i := range sel {

				if groups[i] {
					// If we encounter a new group, and we haven't found any non-nulls for the
					// current group, the output for this group should be null.
					if !a.scratch.foundNonNullForCurrentGroup {
						a.nulls.SetNull(a.curIdx)
					} else {
						a.scratch.vec[a.curIdx] = a.scratch.curSum.Div(int64(a.scratch.curCount))
					}
					a.curIdx++
					a.scratch.curSum = zeroIntervalValue
					a.scratch.curCount = 0

				}

				var isNull bool
				isNull = false
				if !isNull {
					a.scratch.curSum = a.scratch.curSum.Add(col[i])
					a.scratch.curCount++
					a.scratch.foundNonNullForCurrentGroup = true
				}
			}
		}
	}
}

func (a *avgIntervalOrderedAgg) Flush(outputIdx int) {
	// The aggregation is finished. Flush the last value. If we haven't found
	// any non-nulls for this group so far, the output for this group should be
	// NULL.
	// Go around "argument overwritten before first use" linter error.
	_ = outputIdx
	outputIdx = a.curIdx
	a.curIdx++
	if !a.scratch.foundNonNullForCurrentGroup {
		a.nulls.SetNull(outputIdx)
	} else {
		a.scratch.vec[outputIdx] = a.scratch.curSum.Div(int64(a.scratch.curCount))
	}
}

type avgIntervalOrderedAggAlloc struct {
	aggAllocBase
	aggFuncs []avgIntervalOrderedAgg
}

var _ aggregateFuncAlloc = &avgIntervalOrderedAggAlloc{}

const sizeOfAvgIntervalOrderedAgg = int64(unsafe.Sizeof(avgIntervalOrderedAgg{}))
const avgIntervalOrderedAggSliceOverhead = int64(unsafe.Sizeof([]avgIntervalOrderedAgg{}))

func (a *avgIntervalOrderedAggAlloc) newAggFunc() aggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(avgIntervalOrderedAggSliceOverhead + sizeOfAvgIntervalOrderedAgg*a.allocSize)
		a.aggFuncs = make([]avgIntervalOrderedAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	a.aggFuncs = a.aggFuncs[1:]
	return f
}