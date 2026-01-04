"use client";

import type { User } from "@/interface";
import { clearToken, getToken } from "@/utils/auth";
import Link from "next/link";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { useEffect, useRef, useState } from "react";

export function NavLinks() {
  const pathname = usePathname();
  const router = useRouter();
  const searchParams = useSearchParams();
  const q = searchParams.get("query");

  const [query, setQuery] = useState<string>(q ? q : "");
  const [user, setUser] = useState<User | null>(null);
  const [showDropdown, setShowDropdown] = useState(false);
  const [sidebarOpen, setSidebarOpen] = useState(false);
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
    <>
      {/* Sidebar */}
      <div
        className="position-fixed top-0 start-0 bg-dark vh-100"
        style={{
          width: "250px",
          zIndex: 1040,
          transition: "transform 0.3s ease-in-out",
          transform: sidebarOpen ? "translateX(0)" : "translateX(-100%)",
        }}
      >
        <div className="d-flex flex-column h-100 p-3">
          <h5 className="text-light mb-4">Menu</h5>
          <ul className="nav flex-column">
            <li className="nav-item mb-2">
              <Link
                className={`nav-link text-light ${
                  pathname === "/" ? "active bg-primary rounded" : ""
                }`}
                href="/"
                onClick={() => setSidebarOpen(false)}
              >
                <i className="bi bi-house me-2"></i>
                Home
              </Link>
            </li>
            <li className="nav-item mb-2">
              <Link
                className={`nav-link text-light ${
                  pathname === "/watch/list" ? "active bg-primary rounded" : ""
                }`}
                href="/watch/list"
                onClick={() => setSidebarOpen(false)}
              >
                <i className="bi bi-list-ul me-2"></i>
                Watch List
              </Link>
            </li>
            <li className="nav-item mb-2">
              <Link
                className={`nav-link text-light ${
                  pathname === "/watch/calendar"
                    ? "active bg-primary rounded"
                    : ""
                }`}
                href="/watch/calendar"
                onClick={() => setSidebarOpen(false)}
              >
                <i className="bi bi-calendar me-2"></i>
                Watch Calendar
              </Link>
            </li>
          </ul>
        </div>
      </div>

      {/* Overlay */}
      {sidebarOpen && (
        <div
          className="position-fixed top-0 start-0 w-100 h-100 bg-dark bg-opacity-50"
          style={{ zIndex: 1039 }}
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Top Navbar */}
      <nav className="navbar navbar-dark bg-dark">
        <div className="container-fluid">
          <div className="d-flex align-items-center">
            <button
              className="btn btn-link text-light p-0 me-3"
              onClick={() => setSidebarOpen(!sidebarOpen)}
              style={{ fontSize: "1.5rem" }}
            >
              <i className="bi bi-list"></i>
            </button>
            <Link className="navbar-brand" href="/">
              Filmogophery
            </Link>
          </div>

          <form
            className="d-flex position-absolute start-50 translate-middle-x w-25"
            role="search"
            onSubmit={handleSubmit}
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

          <div>
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
      </nav>
    </>
  );
}
