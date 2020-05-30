package sample

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPP_Menu_Paths(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	handle := appRoutes()

	rr := httptest.NewRecorder()
	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Not OK: %v", rr.Code)
	}

	expected := `<h1>MasterPage</h1><p>This is the Home Page</p><p>You are Home!</p> 
	<aside>
		
		<p>
			General
		</p>
		<ul>
			
			<li id="">
				<a  href="/interface">
					Interface
				</a>
				
			</li>
			
		</ul>
		
		<p>
			Stock
		</p>
		<ul>
			
			<li id="">
				<a  href="/stock/parts">
					Parts
				</a>
				
			</li>
			
			<li id="">
				<a  href="/stock/services">
					Services
				</a>
				
			</li>
			
		</ul>
		
	</aside>
	 <span>Footer</span>`
	if rr.Body.String() != expected {
		t.Errorf("unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
