package externalmocks

//go:generate  mockery -output . -outpkg externalmocks -case underscore -dir ../../vendor/gitlab.com/swarmfund/go/doorman -name Doorman
