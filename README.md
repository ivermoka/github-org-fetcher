# github-org-fetcher

Go program that fetches an org's repos from GitHub and grabs name, archived status and URL. 



## Installation and use

1. **Clone the Repository:**
   ```bash
   git clone github.com/ivermoka/github-org-fetcher
  
2. **Navigate to the Repository:**
   ```bash
    cd github-org-fetcher
3. **Run program (add your own org in <org>)**
   ```bash
   go build -o github-org-fetcher github-org-fetcher.go && ./github-org-fetcher <org>

A reposities.json file should now have appeared, containing info on the different org repos. 

**PS: using this program with a large org may result in long load time.**
