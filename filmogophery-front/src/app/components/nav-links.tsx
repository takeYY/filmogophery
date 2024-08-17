"use client";

import { usePathname } from "next/navigation";
import Link from "next/link";
import { useState } from "react";
import { useSearchParams } from "next/navigation";

export function NavLinks() {
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const q = searchParams.get("query");

  const [query, setQuery] = useState<string>(q ? q : "");

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

        <form className="d-flex" role="search">
          <input
            className="form-control me-2 text-light bg-dark"
            type="search"
            placeholder="Search"
            aria-label="Search"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
          <Link
            className="btn btn-outline-primary"
            href={`/search/movie?query=${query}`}
          >
            Search
          </Link>
        </form>
      </div>
    </nav>
  );
}
