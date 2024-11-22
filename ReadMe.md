MIT License

Copyright (c) 2024 Jonas Lambert

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

***

# Fenix Inception

## GuiExecutionServer
GuiExecutionWorker has the responsibility initiate an execution triggered by TesterGui or SuiteScheduler. After storing each TestCase in the database then ExecutionServer is informed over gRPC. The second responsibility is to feed execution statuses, received via Postgres Broadcast System, towards TesterGui using PubSub. 

![Fenix Inception - GuiExecutionServer](./Documentation/FenixInception-Overview-NonDetailed-GuiExecutionServer.png "Fenix Inception - GuiExecutionServer")

The following environment variable is needed for GuiExecutionServer to be able to run.

| Environment variable                        | Example value                                                     | comment                                                                        |
|---------------------------------------------|-------------------------------------------------------------------|--------------------------------------------------------------------------------|
| AuthClientId                                | 46368345345345-au53543bleflkfs03423dfs.apps.googleusercontent.com |                                                                                |
| AuthClientSecret                            | UYGJIU-KHskjshd7HDK7sdfsdf                                        |                                                                                |
| DB_HOST                                     | 127.0.0.1                                                         |                                                                                |
| DB_NAME                                     | fenix_gcp_database                                                |                                                                                |
| DB_PASS                                     | database password                                                 |                                                                                |
| DB_POOL_MAX_CONNECTIONS                     | 4                                                                 | The number of connections towards the database that the ExectionServer can use |
| DB_PORT                                     | 5432                                                              |                                                                                |
| DB_SCHEMA                                   | Not used                                                          |                                                                                |
| DB_USER                                     | postgres                                                          |                                                                                |
| ExecutionLocationForFenixExecutionServer    | GCP                                                               |                                                                                |
| ExecutionLocationForFenixGuiExecutionServer | GCP                                                               |                                                                                |
| FenixExecutionServerAddress                 | fenixexecutionserver-must-be-logged-in-ffafweeerg-lz.a.run.app    |                                                                                |
| FenixExecutionServerPort                    | 443                                                               |                                                                                |
| FenixGuiExecutionServerAddress              | 127.0.0.1                                                         |                                                                                |
| FenixGuiExecutionServerPort                 | 6669                                                              |                                                                                |
| GcpProject                                  | mycloud-run-project                                               |                                                                                |
| LocalServiceAccountPath                     | #                                                                 |                                                                                |
| LogAllSQLs                                  | true                                                              |                                                                                |
| LoggingLevel                                | DebugLevel                                                        |                                                                                |
| TestExecutionStatusPubSubTopicBase          | TestExecutionStatusTopic                                          |                                                                                |
| TestExecutionStatusPubSubTopicSchema        | Fenix-ExecutionStatusUpdateMessageSchema                          |                                                                                |
| UsePubSubWhenSendingExecutionStatus         | true                                                              |                                                                                |

