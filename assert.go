package main

func assert(cond bool, fallback func()) {
	if !cond {
		fallback()
	}
}
