GO=go
APP=TextMiningApp
CMP=TextMiningCompiler

all: app compiler


app:
	${GO} build -o ${APP} src/${APP}/main.go


compiler:
	${GO} build -o ${CMP} src/${CMP}/main.go

test:
	cd src/tests && ${GO} test -bench=.


.PHONY: clean
clean:
	${RM} ${APP}
	${RM} ${CMP}
