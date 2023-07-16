![GitHub Workflow Status](https://github.com/waxdred/GoHotReload/actions/workflows/go.yml/badge.svg)
![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)
# GoHotReload

GoHotReload is a developer tool that automatically reloads any program whenever a file is saved in the repository. It aims to enhance the development experience by eliminating the need for manual restarts during code iterations.

## Features
Automatic reloading: GoHotReload monitors the repository for any changes in the codebase and triggers an automatic reload of the program when a file is saved.

## Table of Contents
- [Prerequisite](#Prerequisite)
- [Configuring](#Configuring)
- [Running](#Running)
- [Structure](#Structure)
- [Monitoring](#Monitoring)
- [Signals](#Signals)

## Prerequisite
Make sure you have Go installed on your machine. You can download it from the official Go website: https://golang.org/dl
Before using GoHotReload, you need to configure the config.json file. This file contains the information about the programs you want to monitor and restart. Each program is defined by the following properties:
- ```path```: The path to the directory where the program is located. Default value is ./.
- ```executable```: The name of the executable file of the program.
- ```extension```: The file extension of the files to monitor. Default value is .go.
- ```cmd```: The command to execute the program.
- ```interval```: The interval in seconds at which the program should be checked for changes. Default value is 4.

## Configuring 
You config multi config.yml files is used to define the programs you want to monitor and restart. Here's an example of how to fill the config.json file:
```yml
configs:
  - name: config1
    cmd:
      - command1 for run your program
      - command2 for run your program
    executable: name of your executable
    extension: .extension
    path: "~/path1"
  
  - name: config2
    cmd:
      - command1 for run your program
      - command2 for run your program
    executable: name of your executable
    extension: .extension
    path: "~/path2"
  
  - name: config3
    cmd:
      - command1 for run your program
    executable: name of your executable
    extension: .extension
    path: "~/path3"
```


![](https://i.imgur.com/Aln3sqf.png)

In this example, we have two programs to monitor. The first program is located in /path/to/program1 directory, has an executable file named program1, and the command to execute it is go run main.go. The second program is located in /path/to/program2 directory, has an executable file named program2, and the command to execute it is python main.py.
Make sure to provide the correct paths, executable names, file extensions, and commands for your programs.

![](https://i.imgur.com/XAz6AJ5.png)

## Running
To run the GoHotReload program, execute the main.go file. The program will read the config.json file, parse the configuration, and start monitoring the specified programs. If there are any errors during the parsing or execution, they will be displayed in the console.
```shell
make
```
you can also build for add on your PATH 
```shell
make build
```
- add on your PATH
- Execute program
```shell
gohot
```
![](https://i.imgur.com/Cy7hLGC.gif)

## Structure
The GoHotReload program consists of the following files:
- ```main.go```: The entry point of the program. It initializes the application and starts the monitoring process.
- ```models/app.go```: Contains the App struct and its methods for parsing the configuration, starting the monitoring process, and handling signals.
- ```models/program.go```: Contains the Program struct and its methods for checking the program path, parsing the configuration, executing the program, and handling the process.
- ```models/handler.go```: Contains the functions for handling signals, checking for file updates, and executing commands.
- ```models/utils.go```: Contains utility functions for printing information and executing system commands.

## Monitoring
The GoHotReload program continuously monitors the specified programs at the specified intervals. It checks for any changes in the files with the specified extension in the program directory. If a file is updated, the program is restarted by executing the specified command.
The program status is displayed in a box format, showing the handler number, status, check status, process status, restart status, executable name, path, command, file extension, process ID (PID), TTY, and memory usage.

## Signals
GoHotReload handles the SIGINT, SIGTERM, and SIGKILL signals to gracefully stop the program. When a signal is received, the program prints a message and exits. It also closes any open file descriptors.

## Conclusion
GoHotReload is a useful tool for monitoring and automatically restarting programs when a file is saved. By configuring the config.json file and running the program, you can easily track and manage multiple programs during development.

## Contributing
Contributions to this project are welcome. If you'd like to contribute, please fork the repository and make your changes. Then, open a pull request and I'll review your changes.

## License
This project is licensed under the MIT License.

