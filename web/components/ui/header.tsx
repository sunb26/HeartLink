"use client";

import { SignedIn, SignedOut, SignInButton, UserButton, useUser } from "@clerk/nextjs";
import { useRouter, usePathname } from "next/navigation";

export function Header() {
  const router = useRouter();
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const pathName = usePathname();

  const { user } = useUser();
  const physicianId = user?.id;
  
  const handleNavigation = (sectionId: string) => {
    router.push(`${sectionId}`);
  };

  return (
    <header className="text-black p-4 text-center top-0 left-0 right-0 border-b-2">
      <nav className="container mx-auto flex justify-between items-center font-[Syne]">
        <div className="font-bold cursor-pointer">
          {/* biome-ignore lint/a11y/useValidAnchor: */}
          <a onClick={() => handleNavigation("/")}>
            <h1 className="text-2xl font-[Syne]">HeartLink</h1>
          </a>
        </div>
        <ul className="flex gap-6 items-center">
          <li>
            <a
              className="hover:underline cursor-pointer text-dark-blue"
              // biome-ignore lint/a11y/useValidAnchor:
              onClick={() => handleNavigation("/")}
            >
              Home
            </a>
          </li>
          <li>
            <a
              className="hover:underline cursor-pointer text-dark-blue"
              // biome-ignore lint/a11y/useValidAnchor:
              onClick={() => handleNavigation("/#about")}
            >
              About
            </a>
          </li>
          <SignedOut>
            <li>
              <SignInButton>
                <button
                  type="button"
                  className="text-white bg-gradient-to-r from-red-400 via-red-500 to-red-600 hover:bg-gradient-to-br focus:ring-4 focus:outline-none focus:ring-red-300 dark:focus:ring-red-800 font-medium rounded-lg text-sm px-3 py-2 text-center me-2 mb-2"
                >
                  Login/Sign Up
                </button>
              </SignInButton>
            </li>
          </SignedOut>
          <SignedIn>
            <li>
              <a
                className="hover:underline cursor-pointer text-dark-blue"
                // biome-ignore lint/a11y/useValidAnchor:
                onClick={() => handleNavigation(`/user/${physicianId}`)}
              >
                My Patients
              </a>
            </li>
            <li>
              <UserButton
                appearance={{
                  elements: {
                    formButtonPrimary:
                      "border-2 border-red-500 bg-white text-red-500 hover:bg-red-500 hover:text-white",
                  },
                }}
              />
            </li>
          </SignedIn>
        </ul>
      </nav>
    </header>
  );
}
