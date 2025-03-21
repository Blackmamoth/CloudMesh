"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Checkbox } from "@/components/ui/checkbox";
import {
  FileIcon,
  FolderIcon,
  DownloadIcon,
  EyeIcon,
  MoreHorizontalIcon,
  GridIcon,
  ListIcon,
  SearchIcon,
  FilterIcon,
  ChevronRightIcon,
  SortIcon,
} from "@/components/icons";
import { MainLayout } from "@/components/layouts/main-layout";
import { FaDropbox, FaGoogleDrive } from "react-icons/fa6";
import { IconType } from "react-icons/lib";

export default function FilesPage() {
  const [viewMode, setViewMode] = useState("list");

  return (
    <MainLayout>
      <div className="flex flex-col gap-6">
        <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div className="flex items-center gap-1 text-muted-foreground">
            <span>Home</span>
            <ChevronRightIcon className="h-4 w-4" />
            <span>Documents</span>
            <ChevronRightIcon className="h-4 w-4" />
            <span className="text-foreground font-medium">Projects</span>
          </div>

          <div className="flex flex-wrap gap-2">
            <div className="relative flex-1 min-w-[200px]">
              <SearchIcon className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input placeholder="Search files..." className="pl-8" />
            </div>

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" size="icon">
                  <FilterIcon className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem>All Files</DropdownMenuItem>
                <DropdownMenuItem>Documents</DropdownMenuItem>
                <DropdownMenuItem>Images</DropdownMenuItem>
                <DropdownMenuItem>Videos</DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline" size="icon">
                  <SortIcon className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem>Name (A-Z)</DropdownMenuItem>
                <DropdownMenuItem>Name (Z-A)</DropdownMenuItem>
                <DropdownMenuItem>Date (Newest)</DropdownMenuItem>
                <DropdownMenuItem>Date (Oldest)</DropdownMenuItem>
                <DropdownMenuItem>Size (Largest)</DropdownMenuItem>
                <DropdownMenuItem>Size (Smallest)</DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="outline">
                  <FaGoogleDrive className="mr-2 h-4 w-4" />
                  All Providers
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem>
                  <FaGoogleDrive className="mr-2 h-4 w-4" />
                  Google Drive
                </DropdownMenuItem>
                <DropdownMenuItem>
                  <FaDropbox className="mr-2 h-4 w-4" />
                  Dropbox
                </DropdownMenuItem>
                <DropdownMenuItem>All Providers</DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>

            <div className="flex rounded-md border">
              <Button
                variant={viewMode === "list" ? "default" : "ghost"}
                size="icon"
                className="rounded-none rounded-l-md"
                onClick={() => setViewMode("list")}
              >
                <ListIcon className="h-4 w-4" />
              </Button>
              <Button
                variant={viewMode === "grid" ? "default" : "ghost"}
                size="icon"
                className="rounded-none rounded-r-md"
                onClick={() => setViewMode("grid")}
              >
                <GridIcon className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>

        <Card>
          <CardContent className="p-0">
            {viewMode === "list" ? <FileListView /> : <FileGridView />}
          </CardContent>
        </Card>
      </div>
    </MainLayout>
  );
}

function FileListView() {
  return (
    <div className="relative w-full overflow-auto">
      <table className="w-full caption-bottom text-sm">
        <thead>
          <tr className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted">
            <th className="h-12 w-12 px-4 text-left align-middle font-medium">
              <Checkbox />
            </th>
            <th className="h-12 px-4 text-left align-middle font-medium">
              Name
            </th>
            <th className="h-12 px-4 text-left align-middle font-medium">
              Size
            </th>
            <th className="h-12 px-4 text-left align-middle font-medium">
              Modified
            </th>
            <th className="h-12 px-4 text-left align-middle font-medium">
              Provider
            </th>
            <th className="h-12 px-4 text-left align-middle font-medium">
              Actions
            </th>
          </tr>
        </thead>
        <tbody>
          <FileListRow
            name="Ashpak Veetar - 64 - cloudmesh.docx"
            Icon={FileIcon}
            size="2.4 MB"
            modified="Today, 2:30 PM"
            Provider={FaGoogleDrive}
          />
          <FileListRow
            name="Backend Software Engineer CV - Ashpak Veetar.pdf"
            Icon={FileIcon}
            size="4.8 MB"
            modified="Yesterday"
            Provider={FaDropbox}
          />
          <FileListRow
            name="CloudMesh.docx"
            Icon={FileIcon}
            size="128 MB"
            modified="Aug 12, 2023"
            Provider={FaGoogleDrive}
          />
          <FileListRow
            name="Copy of abc-id.txt"
            Icon={FileIcon}
            size="8.2 MB"
            modified="Aug 10, 2023"
            Provider={FaGoogleDrive}
          />
          <FileListRow
            name="original.png"
            Icon={FileIcon}
            size="30 MB"
            modified="Aug 5, 2023"
            Provider={FaGoogleDrive}
          />
          <FileListRow
            name="template_bullet-1.docx"
            Icon={FileIcon}
            size="5 MB"
            modified="Aug 5, 2023"
            Provider={FaGoogleDrive}
          />
          <FileListRow
            name="vaultix.excalidraw"
            Icon={FileIcon}
            size="1.2 GB"
            modified="Aug 5, 2023"
            Provider={FaGoogleDrive}
          />
        </tbody>
      </table>
    </div>
  );
}

function FileListRow({
  name,
  Icon,
  size,
  modified,
  Provider,
}: {
  name: string;
  Icon: IconType;
  size: string;
  modified: string;
  Provider: IconType;
}) {
  return (
    <tr className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted">
      <td className="p-4 align-middle">
        <Checkbox />
      </td>
      <td className="p-4 align-middle">
        <div className="flex items-center gap-2">
          <Icon className="h-6 w-6" />
          <span>{name}</span>
        </div>
      </td>
      <td className="p-4 align-middle">{size}</td>
      <td className="p-4 align-middle">{modified}</td>
      <td className="p-4 align-middle">
        <Provider className="h-6 w-6" />
      </td>
      <td className="p-4 align-middle">
        <div className="flex gap-2">
          <Button variant="ghost" size="icon" className="h-8 w-8">
            <EyeIcon className="h-4 w-4" />
          </Button>
          <Button variant="ghost" size="icon" className="h-8 w-8">
            <DownloadIcon className="h-4 w-4" />
          </Button>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="h-8 w-8">
                <MoreHorizontalIcon className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem>Rename</DropdownMenuItem>
              <DropdownMenuItem>Move</DropdownMenuItem>
              <DropdownMenuItem>Copy</DropdownMenuItem>
              <DropdownMenuItem className="text-destructive">
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </td>
    </tr>
  );
}

function FileGridView() {
  const files = [
    {
      name: "Project Proposal.docx",
      icon: <FileIcon className="h-10 w-10" />,
      size: "2.4 MB",
      modified: "Today, 2:30 PM",
      provider: <FaGoogleDrive className="h-4 w-4" />,
    },
    {
      name: "Financial Report.xlsx",
      icon: <FileIcon className="h-10 w-10" />,
      size: "4.8 MB",
      modified: "Yesterday",
      provider: <FaDropbox className="h-4 w-4" />,
    },
    {
      name: "Marketing Assets",
      icon: <FolderIcon className="h-10 w-10" />,
      size: "128 MB",
      modified: "Aug 12, 2023",
      provider: <FaGoogleDrive className="h-4 w-4" />,
    },
    {
      name: "Presentation.pptx",
      icon: <FileIcon className="h-10 w-10" />,
      size: "8.2 MB",
      modified: "Aug 10, 2023",
      provider: <FaDropbox className="h-4 w-4" />,
    },
    {
      name: "Product Photos",
      icon: <FolderIcon className="h-10 w-10" />,
      size: "1.2 GB",
      modified: "Aug 5, 2023",
      provider: <FaGoogleDrive className="h-4 w-4" />,
    },
  ];

  return (
    <div className="grid grid-cols-2 gap-4 p-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
      {files.map((file, index) => (
        <div
          key={index}
          className="group relative rounded-lg border p-3 hover:bg-muted/50"
        >
          <div className="flex flex-col items-center gap-2">
            <div className="flex h-20 w-20 items-center justify-center">
              {file.icon}
            </div>
            <div className="text-center">
              <p className="truncate text-sm font-medium">{file.name}</p>
              <p className="text-xs text-muted-foreground">{file.size}</p>
            </div>
          </div>
          <div className="absolute right-2 top-2 opacity-0 transition-opacity group-hover:opacity-100">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" className="h-8 w-8">
                  <MoreHorizontalIcon className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem>Preview</DropdownMenuItem>
                <DropdownMenuItem>Download</DropdownMenuItem>
                <DropdownMenuItem>Rename</DropdownMenuItem>
                <DropdownMenuItem>Move</DropdownMenuItem>
                <DropdownMenuItem className="text-destructive">
                  Delete
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      ))}
    </div>
  );
}
