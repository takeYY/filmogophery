"use client";

import { usePathname, useRouter } from "next/navigation";
import Link from "next/link";
import { useState } from "react";
import { useSearchParams } from "next/navigation";

export function NavLinks() {
  const pathname = usePathname();
  const router = useRouter();
  const searchParams = useSearchParams();
  const q = searchParams.get("query");

  const [query, setQuery] = useState<string>(q ? q : "");
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (query.trim()) {
      router.push(`/search/movie?query=${query}`);
    }
  };

  return (
    <nav className="navbar nav-underline navbar-expand-lg navbar-dark">
      <div className="container-fluid">
        <a className="navbar-brand" href="#">
          Filmogophery
        </a>
        <button
          className="navbar-toggler"
          type="button"
          data-bs-toggle="collapse"
          data-bs-target="#navbarNavAltMarkup"
          aria-controls="navbarNavAltMarkup"
          aria-expanded="false"
          aria-label="Toggle navigation"
        >
          <span className="navbar-toggler-icon"></span>
        </button>

        <div className="collapse navbar-collapse" id="navbarNavAltMarkup">
          <ul className="navbar-nav ml-auto">
            <li className="nav-item">
              <Link
                className={`nav-link ${pathname === "/" ? "active" : ""}`}
                href="/"
              >
                Home
              </Link>
            </li>

            <li className="nav-item">
              <Link
                className={`nav-link ${
                  pathname === "/watch/list" ? "active" : ""
                }`}
                href="/watch/list"
              >
                Watch List
              </Link>
            </li>

            <li className="nav-item">
              <Link
                className={`nav-link ${
                  pathname === "/watch/calender" ? "active" : ""
                }`}
                href="/watch/calendar"
              >
                Watch Calendar
              </Link>
            </li>
          </ul>
        </div>

        <form className="d-flex" role="search" onSubmit={handleSubmit}>
          <input
            className="form-control me-2 text-light bg-dark"
            type="search"
            placeholder="Search"
            aria-label="Search"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
          <button
            className={`btn btn-outline-primary ${
              !query.trim() ? "disabled" : ""
            }`}
            type="submit"
            disabled={!query.trim()}
          >
            Search
          </button>
        </form>
      </div>
    </nav>
  );
}
