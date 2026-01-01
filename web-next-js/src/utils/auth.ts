interface TokenData {
  accessToken: string;
  refreshToken: string;
  tokenType: string;
  expiresIn: number;
  expiresAt: string;
}

export const saveToken = (token: TokenData) => {
  localStorage.setItem("token", JSON.stringify(token));
};

export const getToken = (): TokenData | null => {
  const token = localStorage.getItem("token");
  return token ? JSON.parse(token) : null;
};

export const clearToken = () => {
  localStorage.removeItem("token");
};
