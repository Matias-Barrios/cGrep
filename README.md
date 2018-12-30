### cGrep

This is a very simple clone of grep which aims to add a little bit of concurrency when more than one file is given as input since it will start a new routine for each one given to.

I got the idea from a book I have been reading too ( which by the way had a different approach altogether ) but I did not base the logic on it as I wanted to make it on my own from scratch as I kind of practice in my golang skills.

This is functional ( insofar as the simplicity of it copes with your expectations ), reads files given as parameters or gets the stdin when it detets there is something there. One problem is that the concurrency is present only when it loads a file so you wont get any benefit from feeding this app from the stdin.


### Flags

**-r [pattern]** This flag is mandatory and it specifies the pattern it should be matched against the file contents.
**-i** It reverses the logic of the pattern matching
**-v** Turns verbose input on ( which basically means it prints the file name and line number along with each match)
**-R [path]** When given it will start a recursive search starting from path given prior to spawning each go routine
**-FF [pattern]** Stands for *file filter* and matches against each file given by the recursive search basically discarding files

### TODO

A million things, but basically, as you might know Go's regex are not as powerful as PCRE ones, so **maybe** in the future Ill add the possibility of using another regex engine then the ones coming with the standard Go's library.

### Install

Clone this repo, cd to the cloned folder and type **make build-linux** or **make build-macos** depending on your OS.
Copy the resulting binary to an appropiate location ( For instance, for linux it would be */bin* or */usr/local/bin*)




