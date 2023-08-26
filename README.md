# overseer
EQEmu Management Suite

Contains the following programs:
- bootstrap: Recommended way to start the server, runs shared memory then overseer
- overseer: Runs zones, world, and other programs, keeping them runnning on crash, and oversees their health
- verify: Verifies if overseer is running properly, and other run time checks
- diagnose: Diagnoses issues with the server, and provides a report to suggest fixes
- install: Installs the server, and all required dependencies
- update: Updates the server, and all required dependencies