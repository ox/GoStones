
//GoStones
//-----------------------
//
//GoStones is an attempt to create RubyGems-like plugin support for Google's Go.
//There will be a main Redmine-run repo called GoForge that will use the little
//script that I created back in the addons.of.cc days.
//
//It works as follows:
//
//	A list of svn/git repos or list providers
//		|			|
//		|			|
//		|		  GoGem +- install -- git/svn pull
//		|			|	|				   |
//		|			|	+---+			   |
//		|	add/remove repo	|	   +-~/go/gems/blah
//		|			|		|	   |		|
//		------------+	uninstall  |		|
//							|	   +---into $PATH
//							|	   |	
//							remove-+
//
//
//Development procedure:
//
//	^-Try to read a file with repos
//	^-Add repo to the list
//	-Try to pull random repos
//	-Try to delete pulled repo
//	-Try to get a certain version of the repo
//	//-Put pulled repo into path
//	-Create a small plugin, have Stones build it
//	//-Put pulled plugin into $PATH, have .go scripts be able to import it { make.bash builds and copies everything }
//
//
//proposed git pull procedure (Frederik Deweerdt):
//	import "os"
//	import "log"
//
//	func main() {
//			var args [3]string;
//			args[0] = "git";
//			args[1] = "clone";
//			args[2] = "YOUR GIT URL HERE";
//			var fds []*os.File = new([3]*os.File);
//			fds[0] = os.Stdin;
//			fds[1] = os.Stdout;
//			fds[2] = os.Stderr;
//
//			/* Replace this with git's full path, or use a shell, and then call git in the args */
//			pid, err := os.ForkExec("/opt/local/bin/git", &args, os.Envs, "/tmp", fds);
//			if err != nil {
//					log.Exit(err)
//			}
//
//			os.Wait(pid, 0);
//	}

package main

import (
	"fmt";
	"io";
	"os";
	"log";
	"bufio";
	"container/vector";
	"flag";
	"path";
	"http";
	"strings";
//	"syscall";
//	"compress/flate";
	)
	
var (
	in *bufio.Reader;
	out *bufio.Writer;
)

var (
	errors vector.Vector;
	repositories vector.StringVector; 
	ik int; 
)

type GemSource struct {
	short_names []string;
	url string;
	}
	
	func NewGemSource(short_names []string, url string) *GemSource {
		return &GemSource{short_names, url}
	}
	
	func (g *GemSource) String() string {
		return g.url
	}
	

func br() {
	fmt.Printf("\n");
}

func git_from_net(url string) {
	var args [3]string;
	args[0] = "git";
	args[1] = "clone";
	args[2] = url;
	var fds []*os.File = new([3]*os.File);
	fds[0] = os.Stdin;
	fds[1] = os.Stdout;
	fds[2] = os.Stderr;

	/* Replace this with git's full path, or use a shell, and then call git in the args */
	pid, err := os.ForkExec("/usr/local/git/bin/git", &args, os.Envs, "./", fds);
	if err != nil {
			log.Exit(err)
	}

	os.Wait(pid, 0);
}


	

func main() {
	errors := vector.New(0);
	
	//fill repositories StringVector
	repositories := vector.New(0); //make a new vec
	file, err := os.Open("./list.txt", os.O_RDONLY, 0755); 
	if file == nil { errors.Push(err); }
	
	in = bufio.NewReader( file );
	ik := 0;
	
	for {
		dat2, err := in.ReadSlice('\n');
		if err != nil || string(dat2[0:len(dat2)-1]) == "" { errors.Push(err); break; }
		
		str := strings.Split(string(dat2), " ", -1);
		gem := NewGemSource(str[1:len(str)], str[0]);
		repositories.Push( gem );
		
		ik++;
		}
	
	
	var get_go_gems_version *bool = flag.Bool("version", false, "Show the version number");
	var show_errors *bool = flag.Bool("show_errors", false, "Show errors");
	var list *bool = flag.Bool("list", false, "List all available repositories");
	var list_full *bool = flag.Bool("list-all", false, "list repositories with their associated aliases");
	var add_repo *string = flag.String("add", "", "Add a repo to the list of known repositories <-add=\"url alias1 alias2\">");
	var remove_repository *string = flag.String("remove-repository", "", "Remove selected repository and it's associated aliases by url or alias");
	var install *string = flag.String("install", "", "Install GoGem");
	flag.Parse();
	
	//get the version number
	if *get_go_gems_version {
		go_gems_version, err := io.ReadFile("./go_gems_version");
		
		if go_gems_version == nil { errors.Push(err) }
		fmt.Println("Welcome to GoGems version ", string(go_gems_version)); 
		
		br();
	}

	//list repositories
	if *list {
		fmt.Println("listing",ik,"repos:");
		for i := 0; i < ik; i++ {
			fmt.Printf("%v\n", repositories.At(i).(*GemSource).url );
			}
			br();
	}
	
	if *list_full {
		for i := 0; i < ik; i++ {
			fmt.Printf("%v\n----------------------\n    ", repositories.At(i).(*GemSource) );
			for _, v := range repositories.At(i).(*GemSource).short_names {
				fmt.Printf("%v\n    ",v);
			}
			br();
		}
	}
	

	
	if *remove_repository != "" {
		
		fmt.Println("looking for", *remove_repository);
		
		l := 0;
		
		for ; l < ik; l++ {
			gem := repositories.At(l).(*GemSource);
			
			if gem.url == *remove_repository {
				fmt.Print("found! removing ", gem.url, " at ", l,"\n");
				repositories.Delete(l);
				break;
			}
			
			for _, v := range gem.short_names {
				if v == *remove_repository || v == *remove_repository + "\n" {
					fmt.Print("found! removing ", v, " at ", l,"\n");
					repositories.Delete(l);
					l = ik;
					break;
				}
			}
		}
				
		if repositories.Len() == ik { //nothing was removed
			//our fall through
			fmt.Println("No such alias or url found in the list, check spelling or url,", repositories.Len());
			os.Exit(1);
		}
				
		file, err := os.Open("./list.txt", os.O_RDWR, 0755); 
		if file == nil { errors.Push(err); }
		
		out = bufio.NewWriter( file );
				
		errlol := io.WriteFile(file.Name(), []byte{}, 0755);
		if errlol != nil { 
			fmt.Print(errlol); 
			errors.Push(errlol); 
			os.Exit(1); 
		}
		
		for i := 0; i < repositories.Len(); i++ {
			gem := repositories.At(i).(*GemSource);
			io.WriteString(out, gem.url);
			fmt.Print(gem);
			
			for k := 0; k < len(gem.short_names); k++ {
				io.WriteString(out, " " + gem.short_names[k]);
				fmt.Print(" " + gem.short_names[k]);
			}
			//io.WriteString(out, "\n" );
		}
				
		out.Flush();	
		file.Close();
	}
	
	if *add_repo != "" {
		str := strings.Split(*add_repo, " ", -1);
		
		var gem *GemSource;
		
		if len(str) == 1 { 
			_, short_name := path.Split(str[0]);
			short_name = strings.Split(short_name, ".",-1)[0];
			gem = NewGemSource([]string{short_name + "\n"}, str[0]);
			fmt.Println("no alias provided, making it", short_name);
		} else {
			gem = NewGemSource(str[1:len(str)], str[0]);
		}
		
		fmt.Println("adding:", gem);
		
		repositories.Push(gem);
		
		file, err := os.Open("./list.txt", os.O_RDWR, 0755); 
		if file == nil { errors.Push(err); }
		
		out = bufio.NewWriter( file );
		
		for i := 0; i < repositories.Len(); i++ {
			io.WriteString(out, repositories.At(i).(*GemSource).url);
			for k := 0; k < len(repositories.At(i).(*GemSource).short_names); k++ {
				io.WriteString(out, " " + repositories.At(i).(*GemSource).short_names[k]);
			}
			//io.WriteString(out, "\n" );
		}
		out.Flush();	
		file.Close();
	}
	
	if *install != "" {
		file_name := *install;
		
		//search the repos for the right file
		for i := 0; i < ik; i++ {
			gem := repositories.At(i).(*GemSource);
			_, pre_tested_name := path.Split( gem.short_names[0] );
			tested_name := strings.Split( pre_tested_name, ".", -1)[0];
			
			for _, val := range gem.short_names {
			
				
			
				if tested_name == file_name || val == file_name || val == file_name+"\n"{
					str := strings.Split(gem.url, ":", -1);
					
					fmt.Println(str[0]);
					
					switch str[0] {
						case "http":
							fmt.Println("Pulling from the net...");
						
							response, _, err := http.Get( gem.url );
							if err != nil { fmt.Println( err ); os.Exit(1); }
							
							var nr int;
							const buf_size = 0x10;
							buf := make ([]byte, buf_size);
							nr, _ = response.Body.Read (buf);

							if nr >= buf_size { panic ("Buffer overrun") }
							
							errorr := io.WriteFile("./"+pre_tested_name, buf, 0755);
							if errorr != nil { fmt.Println(errorr); os.Exit(1); }
							
							buf = nil;
							response.Body.Close ();
							break;
						
						case "git":
							fmt.Println("git'n it from the net...");
							
							git_from_net( string(gem.url) );
							
							break;
						}
					
					fmt.Println("passed retrieving file...");				

					break;
				}
			}
		}
	}
	
	
	if *show_errors {
		fmt.Println(errors);
	}
	
	file.Close();
	
}