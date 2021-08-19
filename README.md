How to trace even with inconsistent Correlation IDs

<img src="https://s.natalian.org/2021-08-19/or.png">

	fields @timestamp, @message, `fields.CORRELATION-ID`, `fields.traceID`| 
	filter `fields.CORRELATION-ID` = "f144b902-33d3-4673-8d27-94c89b766864" or `fields.traceID` = "f144b902-33d3-4673-8d27-94c89b766864"
