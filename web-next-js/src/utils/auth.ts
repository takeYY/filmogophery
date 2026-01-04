interface TokenData {
  accessToken: string;
  refreshToken: string;
  tokenType: string;
  expiresIn: number;
  expiresAt: string;
}

export const saveToken = (token: TokenData) => {
  if (typeof window !== "undefined") {
    localStorage.setItem("token", JSON.stringify(token));
  }
};

export const getToken = (): TokenData | null => {
  if (typeof window === "undefined") return null;
  const token = localStorage.getItem("token");
  return token ? JSON.parse(token) : null;
};

export const clearToken = () => {
  if (typeof window !== "undefined") {
    localStorage.removeItem("token");
  }
};
