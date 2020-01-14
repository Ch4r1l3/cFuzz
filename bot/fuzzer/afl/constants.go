package main

import (
	"regexp"
)

const (
	CORE_PATTERN_FILE_PATH = "/proc/sys/kernel/core_pattern"
	INPUT_DIR_FLAG         = "-i"
	RESUME_FLAG            = "-i-"
	OUTPUT_DIR_FLAG        = "-o"
	MEMORY_LIMIT           = "MEMORY_LIMIT"
	MEMORY_LIMIT_FLAG      = "-m"
	TIMEOUT_LIMIT          = "TIMEOUT_LIMIT"
	TIMEOUT_LIMIT_FLAG     = "-t"
	AFL_PATH               = "/afl/afl-2.52b/afl-fuzz"
	PROGRAM_ARG            = "PROGRAM_ARG"
	TEMP_DIR               = "/tmp"
	AFL_CHECK_TICK_TIME    = 10

	FUZZER_STATS_FILE = "fuzzer_stats"
	CRASH_PATH        = "crashes"
	CRASH_STORE_PATH  = "crashes_store"
)

var CORPUS_CRASH_REGEX = regexp.MustCompile("program crashed with one of the test cases provided")
var STARTUP_CRASH_REGEX = regexp.MustCompile("target binary (crashed|terminated)")
var NO_INSTRUMENTATION_REGEX = regexp.MustCompile("PROGRAM ABORT :.*No instrumentation detected")
var PROGRAM_ABORT_REGEX = regexp.MustCompile("PROGRAM ABORT :.*")
var SANITIZER_START_REGEX = regexp.MustCompile(".*ERROR: [A-z]+Sanitizer:.*")
var AFL_CHECK_REGEX = []*regexp.Regexp{CORPUS_CRASH_REGEX, STARTUP_CRASH_REGEX, NO_INSTRUMENTATION_REGEX, SANITIZER_START_REGEX, PROGRAM_ABORT_REGEX}
var COLOR_REGEX = regexp.MustCompile("\x1B\\[([0-9]{1,2}(;[0-9]{1,2})?)?[mGK]")
