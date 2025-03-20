"use client";

import React, { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { ThemeToggle } from "@/components/theme-toggle";
import {
  CloudIcon,
  MenuIcon,
  SearchIcon,
  BellIcon,
  HomeIcon,
  FolderIcon,
  ArrowRightLeftIcon,
  CloudOffIcon,
  SettingsIcon,
  HelpCircleIcon,
  LogOutIcon,
  UserIcon,
} from "@/components/icons";

export function MainLayout({ children }: { children: React.ReactNode }) {
  const [sidebarOpen, setSidebarOpen] = useState(true);

  return (
    <div className="flex min-h-screen flex-col">
      <header className="sticky top-0 z-50 flex h-16 items-center gap-4 border-b bg-background px-4 md:px-6">
        <Sheet>
          <SheetTrigger asChild>
            <Button variant="outline" size="icon" className="md:hidden">
              <MenuIcon className="h-5 w-5" />
              <span className="sr-only">Toggle Menu</span>
            </Button>
          </SheetTrigger>
          <SheetContent side="left" className="w-72">
            <div className="flex h-full flex-col">
              <div className="flex h-14 items-center border-b px-4">
                <Link
                  href="/dashboard"
                  className="flex items-center gap-2 font-semibold"
                >
                  <CloudIcon className="h-6 w-6" />
                  <span>Cloudmesh</span>
                </Link>
              </div>
              <nav className="grid gap-2 p-4">
                <NavItems />
              </nav>
            </div>
          </SheetContent>
        </Sheet>

        <div className="flex items-center gap-2">
          <Link
            href="/dashboard"
            className="flex items-center gap-2 font-semibold"
          >
            <CloudIcon className="h-6 w-6" />
            <span className="hidden md:inline-block">Cloudmesh</span>
          </Link>
          <Button
            variant="ghost"
            size="icon"
            className="hidden md:flex"
            onClick={() => setSidebarOpen(!sidebarOpen)}
          >
            <MenuIcon className="h-5 w-5" />
            <span className="sr-only">Toggle Sidebar</span>
          </Button>
        </div>

        <div className="relative hidden md:flex md:flex-1">
          <SearchIcon className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Search files..."
            className="w-full max-w-lg pl-8"
          />
        </div>

        <div className="flex items-center gap-2">
          <ThemeToggle />

          <Button variant="ghost" size="icon">
            <BellIcon className="h-5 w-5" />
            <span className="sr-only">Notifications</span>
          </Button>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="rounded-full">
                <Avatar className="h-8 w-8">
                  <AvatarImage
                    src="/placeholder.svg?height=32&width=32"
                    alt="User avatar"
                  />
                  <AvatarFallback>AV</AvatarFallback>
                </Avatar>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <div className="flex items-center gap-2 p-2">
                <div className="flex flex-col space-y-0.5">
                  <span className="text-sm font-medium">Ashpak Veetar</span>
                  <span className="text-xs text-muted-foreground">
                    ashpakv88@gmail.com
                  </span>
                </div>
              </div>
              <DropdownMenuSeparator />
              <DropdownMenuItem asChild>
                <Link href="/profile" className="cursor-pointer">
                  <UserIcon className="mr-2 h-4 w-4" />
                  <span>Profile</span>
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <Link href="/settings" className="cursor-pointer">
                  <SettingsIcon className="mr-2 h-4 w-4" />
                  <span>Settings</span>
                </Link>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem>
                <LogOutIcon className="mr-2 h-4 w-4" />
                <span>Log out</span>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </header>

      <div className="flex flex-1">
        <aside
          className={`${sidebarOpen ? "w-64" : "w-[70px]"} hidden border-r transition-all duration-300 md:block`}
        >
          <nav className="grid gap-2 p-4">
            <NavItems collapsed={!sidebarOpen} />
          </nav>
        </aside>

        <main className="flex-1 p-4 md:p-6">{children}</main>
      </div>
    </div>
  );
}

function NavItems({ collapsed = false }) {
  const pathname = usePathname();

  const navItems = [
    {
      title: "Dashboard",
      href: "/dashboard",
      icon: <HomeIcon className="h-5 w-5" />,
    },
    {
      title: "File Browser",
      href: "/files",
      icon: <FolderIcon className="h-5 w-5" />,
    },
    {
      title: "Transfers",
      href: "/transfers",
      icon: <ArrowRightLeftIcon className="h-5 w-5" />,
    },
    {
      title: "Cloud Connections",
      href: "/connect",
      icon: <CloudOffIcon className="h-5 w-5" />,
    },
    {
      title: "Settings",
      href: "/settings",
      icon: <SettingsIcon className="h-5 w-5" />,
    },
    {
      title: "Help & Support",
      href: "/help",
      icon: <HelpCircleIcon className="h-5 w-5" />,
    },
  ];

  return (
    <>
      {navItems.map((item) => (
        <Button
          key={item.href}
          variant={pathname === item.href ? "default" : "ghost"}
          className={`justify-start ${collapsed ? "px-2" : ""}`}
          asChild
        >
          <Link href={item.href}>
            {item.icon}
            {!collapsed && <span className="ml-2">{item.title}</span>}
          </Link>
        </Button>
      ))}
    </>
  );
}
