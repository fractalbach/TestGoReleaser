# TestGoReleaser

[![Build Status](https://travis-ci.org/fractalbach/TestGoReleaser.svg?branch=master)](https://travis-ci.org/fractalbach/TestGoReleaser)

Experiment with automatic building and deploying and stuff.

asdflkasdfasfd


# Order of Things


1. Human pushes to [TestGoReleaser Github Repo](https://github.com/fractalbach/TestGoReleaser).

2. [Github Webhooks](https://developer.github.com/webhooks/) sends a message to [Travis CI](https://travis-ci.org/).

3. Travis begins.
    - Downloads, builds, and tests TestGoReleaser
    - Only builds once as Linux amd64.
    - Returns success or fail based on the result of tests.

4. Travis Deploy.

    - At this point, Travis might say something like this:   
        ~~~
        Skipping a deployment with the script provider because this is not a tagged commit
        ~~~

    - This step only happens if Travis was prompted by a [Tagged Git Commit](https://git-scm.com/book/en/v2/Git-Basics-Tagging), otherwise Travis stops.

    - Downloads and installs [goreleaser](https://github.com/goreleaser/goreleaser)






# Git Tags

- Check your of your current tags like this:

    ~~~
    git tag
    ~~~



- Making a new *Annotated Tag* looks something this:

    ~~~
    git tag -a v1.4 -m "my version 1.4"
    ~~~

    then push the repo.
