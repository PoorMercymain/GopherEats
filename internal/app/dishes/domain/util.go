package domain

import "time"

var RepeatedAttemptsIntervals = [7]time.Duration{
	getTimeInterval(1),
	getTimeInterval(3),
	getTimeInterval(5),
	getTimeInterval(15),
	getTimeInterval(25),
	getTimeInterval(35),
	getTimeInterval(45),
}

func getTimeInterval(seconds float64) time.Duration {
	value := time.Duration(seconds) * time.Second
	return value
}
