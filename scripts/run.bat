SET target=%1
mkdir bin

cd ..
cd %target%
set GOOS=windows
rsrc -ico overseer.ico -manifest %target%.exe.manifest
go build -buildmode=pie -ldflags="-H=windowsgui -s -w" -o ../bin/%target%.exe .
cd ..
cd bin
%target%.exe
