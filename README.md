# tasksx
A cross-platform system service controller for one-time service and periodical service.

## install

1. Download binary/executable or zip file from the release page for selected platform. 

Run the following command to install it as system-service.

> tasksx install

Note: You need administrator/root priviledge to do this and all the other service related operations(uninstall/reinstall/start/stop, etc).

2. Find out the base path of the service by the following steps:

- First check C:\tasksX in Windows and /tasksX in Linux/MacOS, if the directory exists and (maybe) there are some files in it, then the base path is the directory.
- If not, check the application path(the directory where the tasksx executable is), if there is a file named as testBasePath.txt, then the base path is the application path.

3. Create the task files in the base path.

- for tasks only run once, add or edit the once.json file in base path. The file should in UTF-8 encoding and pure text format as below:

```
[
    {
        "Name": "Once Test 1",
        "Cmd": "dir c:\\ > d:\\tmpx\\testOnceTask1.txt"
    },
    {
        "Name": "Once Test 2",
        "Cmd": "goxc -example basic.gox > d:\\tmpx\\testOnceTask2.txt"
    }
]
```

- for tasks running periodically/repeatly, add or edit the repeat.json file in base path.

```
[
    {
        "Name": "Test Repeat 1",
        "Cmd": "dir c:\\windows > d:\\tmpx\\testRepeatTask1.txt",
        "Start": "2020-04-29 16:00:00",
        "Period": "week"

    },
    {
        "Name": "Test Repeat 2",
        "Cmd": "goxc -example getCurrentTime.gox > d:\\tmpx\\testRepeatTask2.txt",
        "Start": "2020-04-29 16:00:00",
        "Period": "minute"
    },
    {
        "Name": "Test Repeat 3",
        "Cmd": "time /t > d:\\tmpx\\testRepeatTask3.txt",
        "Start": "2020-04-29 16:00:00",
        "Period": "5"
    }
]
```

The Period could be "minute", "hour", "day", "week", or just a number indicates how many minutes. For example, if the Period is "week", the task will repeatly run every week from the Start time. if the Period is "8", the task will repeatly run every 8 minutes. 

4. Restart the service to make the tasks running

run the following command in the console:

> tasksx restart

5. Check the log file for tasks running information

If the base path is C:\tasksX, the log file path will be C:\tasksX\tasksX.log. The content in the log file will be like below:

```
[2020/04/30 09:08:14] ------------------
[2020/04/30 09:08:14] tasksX V0.9a
[2020/04/30 09:08:14] os: windows, basePathG: c:\tasksX, configFileNameG: tasksX.cfg
[2020/04/30 09:08:14] currentPortG: 7489, basePathG: c:\tasksX
[2020/04/30 09:08:14] Service started.
[2020/04/30 09:08:14] Using config file: c:\tasksX\tasksXwin.cfg
[2020/04/30 09:08:14] trying startHttpServer, port: 7489
[2020/04/30 09:08:14] Running repeat task [2] (time /t > d:\tmpx\testRepeatTask3.txt) completed.
[2020/04/30 09:08:14] Running once task [0] (dir c:\ > d:\tmpx\testOnceTask1.txt) completed.
[2020/04/30 09:08:14] Running repeat task [1] (goxc -example getCurrentTime.gox > d:\tmpx\testRepeatTask2.txt) completed.
[2020/04/30 09:08:14] Running once task [1] (goxc -example basic.gox > d:\tmpx\testOnceTask2.txt) completed.
[2020/04/30 09:09:15] Running repeat task [1] (goxc -example getCurrentTime.gox > d:\tmpx\testRepeatTask2.txt) completed.
[2020/04/30 09:10:14] Running repeat task [2] (time /t > d:\tmpx\testRepeatTask3.txt) completed.
[2020/04/30 09:10:15] Running repeat task [1] (goxc -example getCurrentTime.gox > d:\tmpx\testRepeatTask2.txt) completed.
[2020/04/30 09:11:16] Running repeat task [1] (goxc -example getCurrentTime.gox > d:\tmpx\testRepeatTask2.txt) completed.
[2020/04/30 09:12:17] Running repeat task [1] (goxc -example getCurrentTime.gox > d:\tmpx\testRepeatTask2.txt) completed.
[2020/04/30 09:13:17] Running repeat task [1] (goxc -example getCurrentTime.gox > d:\tmpx\testRepeatTask2.txt) completed.
[2020/04/30 09:14:18] Running repeat task [1] (goxc -example getCurrentTime.gox > d:\tmpx\testRepeatTask2.txt) completed.
[2020/04/30 09:15:14] Running repeat task [2] (time /t > d:\tmpx\testRepeatTask3.txt) completed.
[2020/04/30 09:15:18] Running repeat task [1] (goxc -example getCurrentTime.gox > d:\tmpx\testRepeatTask2.txt) completed.
[2020/04/30 09:16:19] Running repeat task [1] (goxc -example getCurrentTime.gox > d:\tmpx\testRepeatTask2.txt) completed.
[2020/04/30 09:17:19] Running repeat task [1] (goxc -example getCurrentTime.gox > d:\tmpx\testRepeatTask2.txt) completed.
[2020/04/30 09:18:19] Running repeat task [1] (goxc -example getCurrentTime.gox > d:\tmpx\testRepeatTask2.txt) completed.
[2020/04/30 09:19:20] Running repeat task [1] (goxc -example getCurrentTime.gox > d:\tmpx\testRepeatTask2.txt) completed.
[2020/04/30 09:20:14] Running repeat task [2] (time /t > d:\tmpx\testRepeatTask3.txt) completed.
[2020/04/30 09:20:20] Running repeat task [1] (goxc -example getCurrentTime.gox > d:\tmpx\testRepeatTask2.txt) completed.
[2020/04/30 09:21:21] Running repeat task [1] (goxc -example getCurrentTime.gox > d:\tmpx\testRepeatTask2.txt) completed.
[2020/04/30 09:22:21] Running repeat task [1] (goxc -example getCurrentTime.gox > d:\tmpx\testRepeatTask2.txt) completed.
```

## Additional Usage Tips

- Uninstall the service

> tasksx uninstall

- Reinstall the service

> tasksx reinstall

- Install the service but not start it

> tasksx installonly

- Manually start an installed service

> tasksx start

- Stop the service

> tasksx stop

## Development/Todo

The program contains a tiny Web server and API server, but currently not used. The default port in 7489.
