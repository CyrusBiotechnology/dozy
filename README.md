# What For?

Dozy is for turning off servers if they do not have any more work to do. It does this by
scanning a given path. If that path resolves to a file, that file's modification 
timestamp is checked against the given (or default) lock validity duration. If the 
timestamp is older than the duration, it is considered stale and disregarded. 

If the path resolves to a folder, the same process as above is run for each file.

Exit code of 0 means there are no valid locks, and all containers have been killed,
otherwise there was an error.

# Usage

Usage help can always be found by running dozy with -h (help). 

    $ dozy -h
    Usage of dozy:
      -duration duration
          duration for which lock files are considered valid (default 10m0s)
      -lock string
          where to look for lockfiles (default "/tmp/lockfiles/")
      -sleep duration
          duration to sleep at the end of script before exit 0


Typical usage looks like:
  
    dozy -duration=30m -lock=/tmp/locksdir/ && shutdown -h now
uptime | sed 's/.*up \([^,]*\), .*/\1/'