# serverlessAndDuplicates

## Task 1:
* Tested with: Python 3.10.6
* See usage instructions `python main.py -h`

**Notes:**
* The "naive" approach for finding duplicates is to read every file and compare its
  contents to all the other files. This will be terribly inefficient because we'll
  have to read all the files times the number of files. (N**2)
* We can avoid this by using a hash table to group duplicate files based on their sha256 signature.
* Because sha256() function can be expensive for very large files,
  we further optimize runtime by first filtering-out unique file sizes.
  If a file size doesn't match any other file size, it cannot be a duplicate
  of anything else thus we don't need to hash it. This optimization assumes that there
  are many big files with different sizes.
* We construct a hash per file in chunks of 8 kb so that we won't run out of memory.


## Task 2:
**Prerequisites:**
1. Terraform +1.5 (- tip #1: use `tfswitch` to get the right version automagically.)
2. Go +1.20
3. [optional] If you want to execute the lambdas locally with `./sam_local_invoke.sh`:
   * SAM CLI, 
   * Docker Desktop (Make sure to export `DOCKER_HOST` environment variable)

**How to deploy this?**
1. Change directories to `./2/terraform`
2. Run `terraform apply`
3. To calculate:
    ```shell
    curl -H "Content-Type: application/json" \
         -X POST $(terraform output -raw apigw_url)/calculate-number-of-occurrences \
         -d '{"word":"example", "character": "e"}'
    ```
4. To get a result:
    ```shell
    curl -H "Content-Type: application/json" \
         -X POST $(terraform output -raw apigw_url)/get-number-of-occurrences \
         -d '{"id":"CHANGE-ME"}'
    ```
    Example:
    ```shell
   $ export BASE_URL=$(terraform output -raw apigw_url | xargs)
   $ curl -H "Content-Type: application/json" \
         -X POST $BASE_URL/calculate-number-of-occurrences \
         -d '{"word":"example", "character": "e"}'
   
    {"result": 2, "id": "b5b1a8c5-2d3d-4931-9707-e37ae43a4522"}
    
    $ curl -H "Content-Type: application/json" \
         -X POST $BASE_URL/get-number-of-occurrences \
         -d '{"id":"b5b1a8c5-2d3d-4931-9707-e37ae43a4522"}'
    
   {"result": 2}
    ```
