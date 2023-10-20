cd ..
mkdir bin
cd overseer
rsrc -ico overseer.ico -manifest overseer.exe.manifest || goto error
go build -buildmode=pie -ldflags="-s -w" -o overseer.exe || goto error
cd ..
copy /y overseer\overseer.exe.manifest bin\overseer.exe.manifest || goto error
move overseer\overseer.exe bin\overseer.exe || goto error
cd bin && overseer.exe c:\games\eq\rebuildeq\rkp.eqg || goto error
exit /b 0

:error
echo Error building %lastdir%: %errorlevel%
exit /b %errorlevel%