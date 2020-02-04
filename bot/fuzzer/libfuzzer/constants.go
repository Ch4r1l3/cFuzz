package main

import (
	"regexp"
)

const (
	MERGE                     = "-merge=1"
	FINAL_STATS               = "-print_final_stats=1"
	MAX_TOTAL_TIME            = "-max_total_time="
	TEMP_DIR                  = "/tmp"
	PROGRAM_ARG               = "PROGRAM_ARG"
	OUTPUT_FLAG               = "-artifact_prefix="
	MAX_MERGE_TIME            = 60
	LIBFUZZER_CHECK_TICK_TIME = 10
)

var STATS_REGEX = regexp.MustCompile("stat::")
