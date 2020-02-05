# cFuzz
Fuzzer Manage System

This is my graduation project

it running fuzzing jobs at scale with Kubernetes

current avaliable fuzzers:
- AFL
- LibFuzzer

you can implement your own fuzzer use go-plugin and upload it to the manage-system

To do the fuzzing jobs, there are several of ways can do it.

1. Use predefined docker image, and upload fuzzer and target
2. Customize your own docker image, you can store fuzzer and target in docker image and use it
