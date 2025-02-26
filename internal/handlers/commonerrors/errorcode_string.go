// Code generated by "stringer -linecomment -type ErrorCode"; DO NOT EDIT.

package commonerrors

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[errUnset-0]
	_ = x[errInternalError-1]
	_ = x[ErrBadValue-2]
	_ = x[ErrFailedToParse-9]
	_ = x[ErrTypeMismatch-14]
	_ = x[ErrAuthenticationFailed-18]
	_ = x[ErrIllegalOperation-20]
	_ = x[ErrNamespaceNotFound-26]
	_ = x[ErrIndexNotFound-27]
	_ = x[ErrUnsuitableValueType-28]
	_ = x[ErrConflictingUpdateOperators-40]
	_ = x[ErrCursorNotFound-43]
	_ = x[ErrNamespaceExists-48]
	_ = x[ErrDollarPrefixedFieldName-52]
	_ = x[ErrInvalidID-53]
	_ = x[ErrEmptyName-56]
	_ = x[ErrCommandNotFound-59]
	_ = x[ErrImmutableField-66]
	_ = x[ErrCannotCreateIndex-67]
	_ = x[ErrInvalidOptions-72]
	_ = x[ErrInvalidNamespace-73]
	_ = x[ErrIndexOptionsConflict-85]
	_ = x[ErrIndexKeySpecsConflict-86]
	_ = x[ErrOperationFailed-96]
	_ = x[ErrDocumentValidationFailure-121]
	_ = x[ErrNotImplemented-238]
	_ = x[ErrDuplicateKey-11000]
	_ = x[ErrStageGroupInvalidFields-15947]
	_ = x[ErrStageGroupID-15948]
	_ = x[ErrStageGroupMissingID-15955]
	_ = x[ErrStageLimitZero-15958]
	_ = x[ErrMatchBadExpression-15959]
	_ = x[ErrProjectBadExpression-15969]
	_ = x[ErrSortBadExpression-15973]
	_ = x[ErrSortBadValue-15974]
	_ = x[ErrSortBadOrder-15975]
	_ = x[ErrSortMissingKey-15976]
	_ = x[ErrStageUnwindWrongType-15981]
	_ = x[ErrPathContainsEmptyElement-15998]
	_ = x[ErrFieldPathInvalidName-16410]
	_ = x[ErrGroupInvalidFieldPath-16872]
	_ = x[ErrGroupUndefinedVariable-17276]
	_ = x[ErrInvalidArg-28667]
	_ = x[ErrSliceFirstArg-28724]
	_ = x[ErrStageUnwindNoPath-28812]
	_ = x[ErrStageUnwindNoPrefix-28818]
	_ = x[ErrProjectionInEx-31253]
	_ = x[ErrProjectionExIn-31254]
	_ = x[ErrAggregatePositionalProject-31324]
	_ = x[ErrWrongPositionalOperatorLocation-31394]
	_ = x[ErrExclusionPositionalProjection-31395]
	_ = x[ErrStageCountNonString-40156]
	_ = x[ErrStageCountNonEmptyString-40157]
	_ = x[ErrStageCountBadPrefix-40158]
	_ = x[ErrStageCountBadValue-40160]
	_ = x[ErrStageGroupUnaryOperator-40237]
	_ = x[ErrStageGroupMultipleAccumulator-40238]
	_ = x[ErrStageGroupInvalidAccumulator-40234]
	_ = x[ErrStageInvalid-40323]
	_ = x[ErrEmptyFieldPath-40352]
	_ = x[ErrInvalidFieldPath-40353]
	_ = x[ErrMissingField-40414]
	_ = x[ErrFailedToParseInput-40415]
	_ = x[ErrCollStatsIsNotFirstStage-40415]
	_ = x[ErrFreeMonitoringDisabled-50840]
	_ = x[ErrValueNegative-51024]
	_ = x[ErrRegexOptions-51075]
	_ = x[ErrRegexMissingParen-51091]
	_ = x[ErrBadRegexOption-51108]
	_ = x[ErrBadPositionalProjection-51246]
	_ = x[ErrElementMismatchPositionalProjection-51247]
	_ = x[ErrEmptySubProject-51270]
	_ = x[ErrEmptyProject-51272]
	_ = x[ErrDuplicateField-4822819]
	_ = x[ErrStageSkipBadValue-5107200]
	_ = x[ErrStageLimitInvalidArg-5107201]
	_ = x[ErrStageCollStatsInvalidArg-5447000]
}

const _ErrorCode_name = "UnsetInternalErrorBadValueFailedToParseTypeMismatchAuthenticationFailedIllegalOperationNamespaceNotFoundIndexNotFoundPathNotViableConflictingUpdateOperatorsCursorNotFoundNamespaceExistsDollarPrefixedFieldNameInvalidIDEmptyFieldNameCommandNotFoundImmutableFieldCannotCreateIndexInvalidOptionsInvalidNamespaceIndexOptionsConflictIndexKeySpecsConflictOperationFailedDocumentValidationFailureNotImplementedLocation11000Location15947Location15948Location15955Location15958Location15959Location15969Location15973Location15974Location15975Location15976Location15981Location15998Location16410Location16872Location17276Location28667Location28724Location28812Location28818Location31253Location31254Location31324Location31394Location31395Location40156Location40157Location40158Location40160Location40234Location40237Location40238Location40323Location40352Location40353Location40414Location40415Location50840Location51024Location51075Location51091Location51108Location51246Location51247Location51270Location51272Location4822819Location5107200Location5107201Location5447000"

var _ErrorCode_map = map[ErrorCode]string{
	0:       _ErrorCode_name[0:5],
	1:       _ErrorCode_name[5:18],
	2:       _ErrorCode_name[18:26],
	9:       _ErrorCode_name[26:39],
	14:      _ErrorCode_name[39:51],
	18:      _ErrorCode_name[51:71],
	20:      _ErrorCode_name[71:87],
	26:      _ErrorCode_name[87:104],
	27:      _ErrorCode_name[104:117],
	28:      _ErrorCode_name[117:130],
	40:      _ErrorCode_name[130:156],
	43:      _ErrorCode_name[156:170],
	48:      _ErrorCode_name[170:185],
	52:      _ErrorCode_name[185:208],
	53:      _ErrorCode_name[208:217],
	56:      _ErrorCode_name[217:231],
	59:      _ErrorCode_name[231:246],
	66:      _ErrorCode_name[246:260],
	67:      _ErrorCode_name[260:277],
	72:      _ErrorCode_name[277:291],
	73:      _ErrorCode_name[291:307],
	85:      _ErrorCode_name[307:327],
	86:      _ErrorCode_name[327:348],
	96:      _ErrorCode_name[348:363],
	121:     _ErrorCode_name[363:388],
	238:     _ErrorCode_name[388:402],
	11000:   _ErrorCode_name[402:415],
	15947:   _ErrorCode_name[415:428],
	15948:   _ErrorCode_name[428:441],
	15955:   _ErrorCode_name[441:454],
	15958:   _ErrorCode_name[454:467],
	15959:   _ErrorCode_name[467:480],
	15969:   _ErrorCode_name[480:493],
	15973:   _ErrorCode_name[493:506],
	15974:   _ErrorCode_name[506:519],
	15975:   _ErrorCode_name[519:532],
	15976:   _ErrorCode_name[532:545],
	15981:   _ErrorCode_name[545:558],
	15998:   _ErrorCode_name[558:571],
	16410:   _ErrorCode_name[571:584],
	16872:   _ErrorCode_name[584:597],
	17276:   _ErrorCode_name[597:610],
	28667:   _ErrorCode_name[610:623],
	28724:   _ErrorCode_name[623:636],
	28812:   _ErrorCode_name[636:649],
	28818:   _ErrorCode_name[649:662],
	31253:   _ErrorCode_name[662:675],
	31254:   _ErrorCode_name[675:688],
	31324:   _ErrorCode_name[688:701],
	31394:   _ErrorCode_name[701:714],
	31395:   _ErrorCode_name[714:727],
	40156:   _ErrorCode_name[727:740],
	40157:   _ErrorCode_name[740:753],
	40158:   _ErrorCode_name[753:766],
	40160:   _ErrorCode_name[766:779],
	40234:   _ErrorCode_name[779:792],
	40237:   _ErrorCode_name[792:805],
	40238:   _ErrorCode_name[805:818],
	40323:   _ErrorCode_name[818:831],
	40352:   _ErrorCode_name[831:844],
	40353:   _ErrorCode_name[844:857],
	40414:   _ErrorCode_name[857:870],
	40415:   _ErrorCode_name[870:883],
	50840:   _ErrorCode_name[883:896],
	51024:   _ErrorCode_name[896:909],
	51075:   _ErrorCode_name[909:922],
	51091:   _ErrorCode_name[922:935],
	51108:   _ErrorCode_name[935:948],
	51246:   _ErrorCode_name[948:961],
	51247:   _ErrorCode_name[961:974],
	51270:   _ErrorCode_name[974:987],
	51272:   _ErrorCode_name[987:1000],
	4822819: _ErrorCode_name[1000:1015],
	5107200: _ErrorCode_name[1015:1030],
	5107201: _ErrorCode_name[1030:1045],
	5447000: _ErrorCode_name[1045:1060],
}

func (i ErrorCode) String() string {
	if str, ok := _ErrorCode_map[i]; ok {
		return str
	}
	return "ErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
}
