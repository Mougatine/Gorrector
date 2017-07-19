GO=go
APP=TextMiningApp
CMP=TextMiningCompiler

all: app compiler


app:
	${GO} build -o ${APP} src/${APP}/main.go


compiler:
	${GO} build -o ${CMP} src/${CMP}/main.go

bench:
	cd src/tests && ${GO} test -bench=.

test:
	cd src/tests && ${GO} test


.PHONY: clean
clean:
	${RM} ${APP}
	${RM} ${CMP}
