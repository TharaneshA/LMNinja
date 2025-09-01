"use client";
import Link from 'next/link';
export default function NotFound() {
  return (
    <div className="flex h-full flex-col items-center justify-center gap-4 text-white">
      <h2 className="text-4xl font-bold">404 - Not Found</h2>
      <p className="text-lg text-gray-400">Could not find the requested resource.</p>
      <Link href="/" className="mt-4 rounded-md bg-[#4A90E2] px-4 py-2 text-white transition-colors hover:bg-[#3A7BC8]">
        Return Home
      </Link>
    </div>
  );
}