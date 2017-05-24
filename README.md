# trueskill

[![Build Status](https://travis-ci.org/mafredri/go-trueskill.svg)](https://travis-ci.org/mafredri/go-trueskill) [![GoDoc](https://godoc.org/mafredri/go-trueskill?status.svg)](https://godoc.org/github.com/mafredri/go-trueskill)

This library implements the [TrueSkill™](http://research.microsoft.com/en-us/projects/trueskill/) ranking system (by Microsoft) in Go.

## TODO

* Refactor the factor graph to remove the need for the distribution bag (collection)
* Support teams and team-based ranking

## Acknowledgements

This implementation is based on [TrueSkill™: A Bayesian Skill Rating System](http://research.microsoft.com/apps/pubs/default.aspx?id=67956) and borrows from the [TrueSkill in F#](http://blogs.technet.com/b/apg/archive/2008/06/16/trueskill-in-f.aspx) test program by Ralf Herbrich. [Computing Your Skill](http://www.moserware.com/2010/03/computing-your-skill.html) by Jeff Moser (and accompanying code) has also been very helpful.
