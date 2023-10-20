cd ..
mkdir bin
cd diagnose
rsrc -ico diagnose.ico -manifest diagnose.exe.manifest
go build -buildmode=pie -ldflags="-s -w" -o diagnose.exe
cd ..
copy /y diagnose\diagnose.exe.manifest bin\diagnose.exe.manifest
move diagnose\diagnose.exe bin\diagnose.exe
cd bin && diagnose.exe c:\games\eq\rebuildeq\rkp.eqg