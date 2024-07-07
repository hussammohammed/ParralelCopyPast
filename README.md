# ParralelCopyPast

# Description:
The Parallel Copy and Paste Tool is a powerful and efficient utility written in Go (Golang) that facilitates the seamless replication of directories, including their nested folders and files. It offers a reliable solution for quickly duplicating directory structures and their contents, saving valuable time and effort in various scenarios, such as software development, data backups, or system administration tasks.

# Key Features:

**Parallel Execution:** Harnessing the power of Go's concurrency capabilities, this tool performs copy and paste operations in parallel, significantly speeding up the replication process for large directories with numerous files and subfolders.

**Recursive Directory Copy:** The tool intelligently replicates entire directory structures, including all subfolders and files contained within. It ensures that the hierarchy is accurately preserved during the copy process, maintaining the integrity of the original directory.

**Efficient File Transfer:** dividing file into chuncks and using Goroutes and channels  to efficiently transfers files, optimizing performance and minimizing the time required to complete the copy operation.

# Usage Example:
```$ parallel-copy-tool -source /path/to/source/directory -destination /path/to/destination/directory```
