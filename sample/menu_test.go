package sample

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMenu_HasCorrectItems(t *testing.T) {
	
}

func TestAPP_Menu_Paths(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	handle := appEpoxy.Router()

	rr := httptest.NewRecorder()
	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Not OK: %v", rr.Code)
	}

	expected := `<h1>MasterPage</h1><p>This is the Home Page</p><p>Welcome</p>
	<aside>
	<p>
			Home
		</p>
		<ul>
			<li><a href="/stock">Stock</a></li>
			<li>
				<a href="/stock/parts">Parts</a>
				<ul>
					<li><a href="/stock/parts/create">Create</a></li>
				</ul>
			</li>
			<li>
			<a href="/stock/services">Services</a>
			<ul>
				<li><a href="/stock/services/create">Create</a></li>
			</ul>
		</li>
		</ul>
		<p>
			Stock.API
		</p>
		<ul>
			<li><a href="/stock">Stock</a></li>
			<li>
				<a href="/stock/parts">Parts</a>
				<ul>
					<li><a href="/stock/parts/create">Create</a></li>
				</ul>
			</li>
			<li>
			<a href="/stock/services">Services</a>
			<ul>
				<li><a href="/stock/services/create">Create</a></li>
			</ul>
		</li>
		</ul>
	</aside>`
	if rr.Body.String() != expected {
		t.Errorf("unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
