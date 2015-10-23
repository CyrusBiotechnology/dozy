# What For?

Dozy is for turning off servers if they do not have any more work to do. It does this by
scanning a given path. If that path resolves to a file, that file's modification
timestamp is checked against the given (or default) lock validity duration. If the 
timestamp is older than the duration, it is considered stale and disregarded.

If the path resolves to a folder, the same process as above is run for each file found
within that folder, recursively.

Exit code of 0 means there are no valid locks, and all containers have been killed,
otherwise there was an error.

# Requirements

`dozy` targets Darwin (OS X) and Linux environments. Other systems are not currently 
supported.

# Installing

Simply download the appropriate "dozy" release for your platform and copy the executable
to someplace on your path. `/usr/local/bin` is often a good spot. 

# Usage

Usage help can always be found by running dozy with -h (help). 

    $ dozy -h
    
    Usage of dozy:
      -duration duration
          duration for which lock files are considered valid (default 10m0s)
      -lock string
          where to look for lockfiles (default "/tmp/lockfiles/")
      -minuptime duration
          will not exit 0 before uptime >= <minuptime>
      -sleep duration
          duration to sleep at the end of script before exit 0


Typical usage (in a [crontab](https://en.wikipedia.org/wiki/Cron)) looks something like:

    # M H DoM Mo DoW    COMMAND
      * * *   *  *      /usr/local/bin/dozy -minuptime=120m -duration=30m -lock=/tmp/locksdir/ && shutdown -h now
