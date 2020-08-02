## Prerequisite

- Go 1.14 or latest

- WSL2, Linux or macOS

- PostgreSQL

- Docker

- Make


## Boilerplate Overview

- Setup and Command
Clone this repository: https://github.com/tunaiku/smb-recruitment-be - Connect to preview 

- Open your terminal and open the boilerplate directory

- Make sure that you are running docker

  - For set up the project environment, please follow a certain step below :

  - Run this command on your terminal in order to configure the project environment
    ```
    > make setupposgres
    > make setupdata
    ```
    You are ready to go

  - For running the end-to-end test case can use this command:
    ```
    > make e2e
    ```
  - For building the project, you can use this command:
    ```
    > make buildapp
    ```
  - For running the project, you can use this command:
    ```
    > make run
    ```

