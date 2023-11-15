# github-org-fetcher

Go program that fetches an org's repos from GitHub and grabs name, archived status and URL. 



## Installation and use

1. **Clone the Repository:**
   ```bash
   git clone github.com/ivermoka/github-org-fetcher
  
2. **Navigate to the Repository:**
   ```bash
    cd github-org-fetcher
3. **Generate and copy GitHub API authentication token (in settings -> developer settings -> personal access tokens -> Token (classic))**.
 Select the options shown on the image below:
   <img width="1216" alt="Screenshot 2023-11-15 at 19 25 06" src="https://github.com/ivermoka/github-org-fetcher/assets/119415554/b6eacb93-e307-43d0-92a9-360e2fd4286d">

   Scroll down and select "Generate token". Click copy. 
5. **Run program (add your own org in <org>)**
   ```bash
   go build -o github-org-fetcher github-org-fetcher.go && ./github-org-fetcher -a <token> <org>

A reposities.json file should now have appeared, containing info on the different org repos. 

**PS: using this program with a large org may result in long load time.**

