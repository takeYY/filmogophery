"use client";

import { usePathname, useRouter } from "next/navigation";
import Link from "next/link";
import { useState, useEffect, useRef } from "react";
import { useSearchParams } from "next/navigation";
import { getToken, clearToken } from "@/utils/auth";
import type { User } from "@/interface";

export function NavLinks() {
  const pathname = usePathname();
  const router = useRouter();
  const searchParams = useSearchParams();
  const q = searchParams.get("query");

  const [query, setQuery] = useState<string>(q ? q : "");
  const [user, setUser] = useState<User | null>(null);
  const [showDropdown, setShowDropdown] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const fetchUser = async () => {
      const token = getToken();
      if (!token) return;

      try {
        const res = await fetch("/api/users/me", {
          headers: {
            Authorization: `${token.tokenType} ${token.accessToken}`,
          },
        });
        if (res.ok) {
          const data = await res.json();
          setUser(data);
        }
      } catch (err) {
        console.error("Failed to fetch user", err);
      }
    };

    fetchUser();
  }, []);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setShowDropdown(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (query.trim()) {
      router.push(`/search/movie?query=${query}`);
    }
  };

  const handleLogout = async () => {
    const token = getToken();
    if (token) {
      await fetch("/api/auth/logout", {
        method: "POST",
        headers: {
          Authorization: `${token.tokenType} ${token.accessToken}`,
        },
      });
    }
    clearToken();
    setUser(null);
    router.push("/login");
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
          <ul className="navbar-nav">
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

          <form
            className="d-flex position-absolute start-50 translate-middle-x"
            role="search"
            onSubmit={handleSubmit}
            style={{ width: "400px" }}
          >
            <input
              className="form-control me-2 text-light bg-dark"
              type="search"
              placeholder="Search"
              aria-label="Search"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
            />
            <button
              className="btn btn-outline-primary"
              type="submit"
              disabled={!query.trim()}
            >
              <i className="bi bi-search"></i>
            </button>
          </form>

          <div className="ms-auto">
            {user && (
              <div className="position-relative" ref={dropdownRef}>
                <button
                  className="btn btn-link text-light p-0"
                  onClick={() => setShowDropdown(!showDropdown)}
                  style={{ fontSize: "1.5rem" }}
                >
                  <i className="bi bi-person-circle"></i>
                </button>
                {showDropdown && (
                  <div
                    className="dropdown-menu dropdown-menu-end show position-absolute"
                    style={{ right: 0, minWidth: "200px" }}
                  >
                    <div className="px-3 py-2 border-bottom">
                      <div className="fw-bold">{user.username}</div>
                      <small className="text-muted">{user.email}</small>
                    </div>
                    <button className="dropdown-item" onClick={handleLogout}>
                      ログアウト
                    </button>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}
