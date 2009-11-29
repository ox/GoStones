/*
GoGems
-----------------------

GoGems is an attempt to create RubyGems-like plugin support for Google's Go.
There will be a main Redmine-run repo called GoForge that will use the little
script that I created back in the addons.of.cc days.

It works as follows:

	A list of svn/git repos or list providers
		|			|
		|			|
		|		  GoGem +- install -- git/svn pull
		|			|	|				   |
		|			|	+---+			   |
		|	add/remove repo	|	   +-~/go/gems/blah
		|			|		|	   |		|
		------------+	uninstall  |		|
							|	   +---into $PATH
							|	   |	
							remove-+


Development procedure:

	^-Try to read a file with repos
	^-Add repo to the list
	-Try to pull random repos
	-Try to delete pulled repo
	-Try to get a certain version of the repo
	-Put pulled repo into path
	-Create a small plugin, have Gems build it
	-Put pulled plugin into $PATH, have .go scripts be able to import it


*/

package main

import (
	"fmt";
	"io";
	"os";
	"bufio";
	"container/vector";
	"flag";
	"path";
	"http";
	"strings";
//	"compress/flate";
	)
	
var (
	in *bufio.Reader;
	out *bufio.Writer;
)

var (
	repositories vector.StringVector; 
	ik int; 
)

type GemSource struct {
	short_name string;
	url string;
	}

	

func main() {
	errors := vector.New(100);
	
	//fill repositories StringVector
	repositories := vector.NewStringVector(0); //make a new vec
	file, err := os.Open("./list.txt", os.O_RDONLY, 0755); 
	if file == nil { errors.Push(err); }
	
	in = bufio.NewReader( file );
	ik := 0;
	
	for {
		dat2, err := in.ReadSlice('\n');
		if err != nil || string(dat2[0:len(dat2)-1]) == "" { errors.Push(err); break; }
		repositories.Push( string(dat2[0:len(dat2)-1]) );
		
		ik++;
		}
	
	
	var get_go_gems_version *bool = flag.Bool("version", false, "Show the version number");
	var show_errors *bool = flag.Bool("show_errors", false, "Show errors");
	var list *bool = flag.Bool("list", false, "List all available repositories");
	var add_repo *string = flag.String("add", "", "Add a repo to the list of known repositories");
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
			fmt.Println( repositories.At(i) );
			}
			br();
	}
	
	if *add_repo != "" {
		repositories.Push(*add_repo);
		
		file, err := os.Open("./list.txt", os.O_RDWR, 0755); 
		if file == nil { errors.Push(err); }
		
		out = bufio.NewWriter( file );
		
		for i := 0; i < ik+1; i++ {
			io.WriteString(out, repositories.At(i)+"\n" );
		}
		out.Flush();	
		file.Close();	
	}
	
	if *install != "" {
		file_name := *install;
		
		//search the repos for the right file
		for i := 0; i < ik; i++ {
			_, pre_tested_name := path.Split( repositories.At(i) );
			tested_name := strings.Split( pre_tested_name, ".", -1)[0];
			
			if tested_name == file_name {
				response, _, err := http.Get( repositories.At(i) );
				if err != nil { fmt.Println( err ); }
				
				var nr int;
				const buf_size = 0x1000;
				buf := make ([]byte, buf_size);
				nr, _ = response.Body.Read (buf);

				if nr >= buf_size { panic ("Buffer overrun") }
				
				errorr := io.WriteFile("./"+pre_tested_name, buf, 0755);
				if errorr != nil { fmt.Println(errorr) }
				
				buf = nil;
				response.Body.Close ();
				
				fmt.Println("passed retrieving file...");
				
				//now lets handle unzipping, so far only gzip and zlib
//				file, _ := os.Open("./"+pre_tested_name, os.O_RDONLY, 0755);
//				reader := bufio.NewReader( file );
//				infl := flate.NewInflater( reader );
//				
//				buf2 := make([]byte, buf_size);
//				deflated, _ := infl.Read(buf2);
//				if deflated >= buf_size { panic("BUFFER OVERRUN!") }
//				
//				errr := io.WriteFile("./"+tested_name, buf2, 0755);
//				if errr != nil { fmt.Println(errr) };
//				buf2 = nil;
//				
//				infl.Close();
				
				
				break;
			}
		}
	}
	
	
	if *show_errors {
		fmt.Println(errors);
	}
}