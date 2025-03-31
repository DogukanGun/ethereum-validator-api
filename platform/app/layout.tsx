import type { Metadata, Viewport } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const viewport: Viewport = {
  themeColor: "#1e293b",
  width: "device-width",
  initialScale: 1.0,
};

export const metadata: Metadata = {
  title: "Ethereum Validator Explorer",
  description: "Explore Ethereum validator performance, sync committee duties, and block rewards across the Beacon Chain.",
  keywords: ["ethereum", "validator", "blockchain", "beacon chain", "staking", "sync committee", "block rewards"],
  authors: [{ name: "Dogukan Gundogan" }],
  openGraph: {
    title: "Ethereum Validator Explorer",
    description: "Explore Ethereum validator performance and rewards",
    type: "website",
    locale: "en_US",
    images: [
      {
        url: "/og-image.png",
        width: 1200,
        height: 630,
        alt: "Ethereum Validator Explorer",
      },
    ],
  },
  twitter: {
    card: "summary_large_image",
    title: "Ethereum Validator Explorer",
    description: "Explore Ethereum validator performance and rewards",
    images: ["/og-image.png"],
  },
  icons: {
    icon: [
      { url: "/favicon.ico" },
      { url: "/favicon-16x16.png", sizes: "16x16", type: "image/png" },
      { url: "/favicon-32x32.png", sizes: "32x32", type: "image/png" },
    ],
    apple: [
      { url: "/apple-touch-icon.png" },
    ],
    other: [
      {
        rel: "mask-icon",
        url: "/safari-pinned-tab.svg",
        color: "#5bbad5",
      },
    ],
  },
  manifest: "/site.webmanifest",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        {children}
      </body>
    </html>
  );
}
