"use client";

import { useAuth } from "@/hooks/useAuth";
import { MyWatchHistory } from "@/interface/index";
import { useEffect, useState } from "react";
import Link from "next/link";
import Image from "next/image";
import { posterUrlPrefix } from "@/constants/poster";
import styles from "./calendar.module.css";

export default function Page() {
  const [loading, setLoading] = useState(true);
  const [watchHistory, setWatchHistory] = useState<MyWatchHistory[]>([]);
  const [currentDate, setCurrentDate] = useState(new Date());

  const token = useAuth();
  const accessToken = token ? token.accessToken : null;

  const headers: HeadersInit = {};
  if (accessToken) {
    headers.Authorization = `Bearer ${accessToken}`;
  }

  useEffect(() => {
    const fetchWatchHistory = async () => {
      setLoading(true);
      try {
        const response = await fetch(`/api/users/me/watch-history?limit=100`, {
          method: "GET",
          headers,
        });
        const data: MyWatchHistory[] = await response.json();
        setWatchHistory(data);
      } catch {
        setWatchHistory([]);
      } finally {
        setLoading(false);
      }
    };
    fetchWatchHistory();
  }, []);

  const year = currentDate.getFullYear();
  const month = currentDate.getMonth();

  const firstDay = new Date(year, month, 1).getDay();
  const daysInMonth = new Date(year, month + 1, 0).getDate();

  const prevMonth = () => setCurrentDate(new Date(year, month - 1, 1));
  const nextMonth = () => setCurrentDate(new Date(year, month + 1, 1));

  const getWatchedMovies = (day: number) => {
    const dateStr = `${year}-${String(month + 1).padStart(2, "0")}-${String(day).padStart(2, "0")}`;
    return watchHistory.filter((wh) => wh.watchedAt.startsWith(dateStr));
  };

  const isToday = (day: number) => {
    const today = new Date();
    return today.getFullYear() === year && today.getMonth() === month && today.getDate() === day;
  };

  const renderCalendar = () => {
    const days = [];
    for (let i = 0; i < firstDay; i++) {
      days.push(<div key={`empty-${i}`} className={`${styles.day} ${styles.dayEmpty}`}></div>);
    }
    for (let day = 1; day <= daysInMonth; day++) {
      const movies = getWatchedMovies(day);
      const today = isToday(day);
      days.push(
        <div key={day} className={`${styles.day} ${today ? styles.dayToday : ""}`}>
          <div className={styles.dayNumber}>{day}</div>
          {movies.length > 0 && <div className={styles.movieCount}>{movies.length}</div>}
          <div className={styles.posterGrid}>
            {movies.slice(0, 4).map((wh) => (
              <Link key={wh.id} href={`/movie/${wh.movie.id}`} className={styles.posterLink}>
                <Image
                  src={posterUrlPrefix + (wh.movie.posterURL || "/Agz71U0wcesx87micVn731Z1vPu.jpg")}
                  alt={wh.movie.title}
                  width={92}
                  height={138}
                  className={styles.poster}
                />
              </Link>
            ))}
          </div>
        </div>
      );
    }
    return days;
  };

  return (
    <div className="container pb-4">
      <h3 className="text-center mb-4 text-light">Watch Calendar</h3>

      {loading ? (
        <div className={styles.loading}>
          <div className={styles.spinner}></div>
        </div>
      ) : (
        <>
          <div className={styles.header}>
            <button className={styles.navButton} onClick={prevMonth}>
              ← 前月
            </button>
            <h2 className={styles.title}>
              {year}年 {month + 1}月
            </h2>
            <button className={styles.navButton} onClick={nextMonth}>
              次月 →
            </button>
          </div>

          <div className={styles.calendar}>
            <div className={styles.weekdays}>
              {["日", "月", "火", "水", "木", "金", "土"].map((day) => (
                <div key={day} className={styles.weekday}>
                  {day}
                </div>
              ))}
            </div>
            <div className={styles.days}>
              {renderCalendar()}
            </div>
          </div>
        </>
      )}
    </div>
  );
}
